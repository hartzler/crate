#define _GNU_SOURCE
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>

#include <linux/limits.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <signal.h>
#include <setjmp.h>
#include <sched.h>
#include <signal.h>

/* All arguments should be above stack, because it grows down */
struct clone_arg {
	/*
	 * Reserve some space for clone() to locate arguments
	 * and retcode in this place
	 */
	char stack[4096] __attribute__ ((aligned(8)));
	char stack_ptr[0];
	jmp_buf *env;
};

static int child_func(void *_arg)
{
	struct clone_arg *arg = (struct clone_arg *)_arg;
	longjmp(*arg->env, 1);
}

#define pr_perror(fmt, ...) fprintf(stderr, "nsenter: " fmt ": %m\n", ##__VA_ARGS__)

// Use raw setns syscall for versions of glibc that don't include it (namely glibc-2.12)
#if __GLIBC__ == 2 && __GLIBC_MINOR__ < 14
#define _GNU_SOURCE
#include "syscall.h"
#ifdef SYS_setns
int setns(int fd, int nstype)
{
	return syscall(SYS_setns, fd, nstype);
}
#endif
#endif

static int clone_parent(jmp_buf * env) __attribute__ ((noinline));
static int clone_parent(jmp_buf * env)
{
	struct clone_arg ca;
	int child;

	ca.env = env;
	child = clone(child_func, ca.stack_ptr, CLONE_PARENT | SIGCHLD, &ca);

	return child;
}

void nsexec()
{
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };
	const int num = sizeof(namespaces) / sizeof(char *);
	jmp_buf env;
	char buf[PATH_MAX], *val;
	int i, tfd, child, len;
	pid_t pid;

	val = getenv("_LIBCONTAINER_INITPID");
	if (val == NULL)
		return;

	pid = atoi(val);
	snprintf(buf, sizeof(buf), "%d", pid);
	if (strcmp(val, buf)) {
		pr_perror("Unable to parse _LIBCONTAINER_INITPID");
		exit(1);
	}

	/* Check that the specified process exists */
	snprintf(buf, PATH_MAX - 1, "/proc/%d/ns", pid);
	tfd = open(buf, O_DIRECTORY | O_RDONLY);
	if (tfd == -1) {
		pr_perror("Failed to open \"%s\"", buf);
		exit(1);
	}

	for (i = 0; i < num; i++) {
		struct stat st;
		int fd;

		/* Symlinks on all namespaces exist for dead processes, but they can't be opened */
		if (fstatat(tfd, namespaces[i], &st, AT_SYMLINK_NOFOLLOW) == -1) {
			// Ignore nonexistent namespaces.
			if (errno == ENOENT)
				continue;
		}

		fd = openat(tfd, namespaces[i], O_RDONLY);
		if (fd == -1) {
			pr_perror("Failed to open ns file %s for ns %s", buf,
				  namespaces[i]);
			exit(1);
		}
		// Set the namespace.
		if (setns(fd, 0) == -1) {
			pr_perror("Failed to setns for %s", namespaces[i]);
			exit(1);
		}
		close(fd);
	}

	if (setjmp(env) == 1) {
		// Finish executing, let the Go runtime take over.
		return;
	}

	child = clone_parent(&env);
	if (child < 0) {
		pr_perror("Unable to fork");
		exit(1);
	}

	len = snprintf(buf, sizeof(buf), "{ \"pid\" : %d }\n", child);

	if (write(3, buf, len) != len) {
		pr_perror("Unable to send a child pid");
		kill(child, SIGKILL);
		exit(1);
	}

	exit(0);
}
