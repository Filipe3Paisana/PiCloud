# PiCloud - Decentralized File Storage Platform

PiCloud is an innovative academic project offering decentralized file storage services, ensuring security, scalability, and reliability. With PiCloud, users can upload, download, and manage files seamlessly while contributing to the network by hosting a node on their devices.

---

## 🚀 Features

- **Decentralized Storage**: Files are fragmented and stored across multiple geographically distributed nodes.
- **User-Friendly Interface**: Intuitive web interface for file uploads, downloads, and management.
- **Advanced Security**: TLS for data in transit, AES for storage, and robust authentication mechanisms.
- **AI Integration**: Predicts node performance and redistributes file fragments for optimal resource usage.
- **Scalability**: Supports horizontal scaling by adding more nodes to the network.
- **Modular Architecture**: Built with Docker and Kubernetes for efficient deployment and maintenance.

---

## 🛠️ Technologies Used

- **Backend**: [Go](https://golang.org/) for API and node services.
- **Frontend**: HTML, CSS, JavaScript for a responsive and intuitive web experience.
- **Containerization**: Docker for isolating components and Kubernetes for orchestration.
- **Security**: Cloudflare for enhanced protection and secure exposure of APIs.
- **AI Models**: HistGradientBoostingRegressor for predicting node performance.
- **Hardware**: Raspberry Pi for local node hosting.

---

## 📚 How It Works

1. **File Upload**: Files are fragmented and metadata is stored in a central database.
2. **File Distribution**: Fragments are distributed across nodes, ensuring redundancy and availability.
3. **File Retrieval**: Fragments are reassembled and delivered securely to the user.
4. **Node Monitoring**: AI models monitor node health and redistribute fragments as needed.

---

## 🔧 Getting Started

### Prerequisites

- Docker and Docker Compose installed.
- Kubernetes cluster (optional for advanced deployment).
- Raspberry Pi or other devices for hosting nodes.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/picloud.git
   cd picloud
