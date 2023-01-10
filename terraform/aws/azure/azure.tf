# Specificare il provider di cloud
provider "azurerm" {
  subscription_id = "SUBSCRIPTION_ID"
  client_id       = "CLIENT_ID"
  client_secret   = "CLIENT_SECRET"
  tenant_id       = "TENANT_ID"
}

# Creare una risorsa di tipo "azurerm_virtual_machine"
resource "azurerm_virtual_machine" "web" {
  # Scegliere un'immagine di sistema operativo
  storage_os_disk {
    image_reference {
      publisher = "Canonical"
      offer     = "UbuntuServer"
      sku       = "20.04-LTS"
      version   = "latest"
    }
  }
  vm_size  = "Standard_B1s"

  # Assegnare un nome di host
  name    = "web"
  location = "westus2"

  # Aprire le porte SSH e HTTPS
  network_interface_ids = [azurerm_network_interface.web.id]

  # Assegnare un indirizzo IP pubblico alla VM
  public_ip_address_id = azurerm_public_ip.web.id
}

# Creare una risorsa di tipo "azurerm_network_interface"
resource "azurerm_network_interface" "web" {
  name                = "web"
  location            = azurerm_virtual_machine.web.location
  resource_group_name = azurerm_virtual_machine.web.resource_group_name

  # Creare una regola di firewall per aprire la porta 22
  security_group = azurerm_network_security_group.web
}

# Creare una risorsa di tipo "azurerm_network_security_group" per aprire le porte SSH e HTTPS
resource "azurerm_network_security_group" "web" {
  name                = "web"
  location            = azurerm_virtual_machine.web.location
  resource_group_name = azurerm_virtual_machine.web.resource_group_name

  # Aprire le porte 22 e 443 per il traffico in entrata
  security_rule {
    name                       = "ssh"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "https"
    priority                   = 1002
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "*"
    destination_address
  }

}


output "public_ip" {
  value = azurerm_public_ip.web.ip_address
}
