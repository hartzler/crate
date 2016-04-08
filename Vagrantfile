# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/wily64"
  config.vm.synced_folder '.', '/vagrant/src/github.com/armada-io/crate'
  config.vm.provision :shell, path: 'script/provision.sh'
  config.vm.provision :shell, path: 'script/mount-cgroups.sh'
  config.vm.provider "virtualbox" do |v|
    v.memory = 1024
    v.cpus = 2
  end
end
