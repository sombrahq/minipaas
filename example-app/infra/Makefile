MPAAS=../../cli/build/minipaas
ANSIBLE_PLAYBOOK=install.yml
INVENTORY=inventory.ini

.PHONY: up ssh destroy rebuild provision

up:
	vagrant up --provider=qemu

certs:
	$(MPAAS) certs server --cn=localhost --verbose

provision:
	ansible-playbook -i inventory.ini install.yml

destroy:
	vagrant destroy -f

ssh:
	vagrant ssh

rebuild: destroy up
