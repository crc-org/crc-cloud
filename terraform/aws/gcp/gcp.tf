# Specificare il provider di cloud
provider "google" {
  credentials = file("credentials.json")
  project     = "my-project"
  region      = "us-central1"
}

# Creare una risorsa di tipo "google_compute_instance"
resource "google_compute_instance" "web" {
  # Scegliere un'immagine di sistema operativo
  machine_type = "f1-micro"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-10"
    }
  }

  # Assegnare un nome di host
  metadata = {
    "name" = "web"
  }

  # Aprire le porte SSH e HTTPS
  network_interface {
    network = "default"

    access_config {
      # Assegnare un indirizzo IP pubblico alla VM
      nat_ip = "auto"
    }

    # Aggiungere regole di firewall per aprire le porte 22 e 443
    allowed_ssh_ports = ["22"]
    allowed_https_ports = ["443"]
  }
}

# Creare una risorsa di tipo "google_compute_firewall" per aprire la porta SSH
resource "google_compute_firewall" "ssh" {
  name        = "ssh"
  description = "Allow SSH access"
  network     = "default"

  # Aprire la porta 22 per il traffico in entrata
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  # Aprire anche la porta 6443
  allow {
    protocol = "tcp"
    ports    = ["6443"]
  }

  # Accettare il traffico in entrata da qualsiasi indirizzo IP
  source_ranges = ["0.0.0.0/0"]
}

# Creare una risorsa di tipo "google_compute_firewall" per aprire la porta HTTPS
resource "google_compute_firewall" "https" {
  name        = "https"
  description = "Allow HTTPS access"
  network     = "default"

  # Aprire la porta 443 per il traffico in entrata
  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  # Accettare il traffico in entrata da qualsiasi indirizzo IP
  source_ranges = ["0.0.0.0/0"]
}

output "public_ip" {
  value = google_compute_instance.web.network_interface.0.access_config.0.nat_ip
}

