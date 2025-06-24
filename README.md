# Bluetooth Advertisement Sniffer

A Go-based Bluetooth Low Energy (BLE) advertisement data sniffer that can capture and display BLE advertisement packets with optional device address filtering.

## Installation

```bash
inputs = {
  bluetooth-sniffer.url = "github:nxm/bluetooth-sniffer";
};

home.packages = [
  inputs.bluetooth-sniffer.packages.${pkgs.system}.default
];
```

## Usage

### Basic scanning (all devices)
```bash
sudo ./bluetooth-sniffer
```

### Filter by device address
```bash
# Filter devices containing "aa:bb" in their address
sudo ./bluetooth-sniffer -addr aa:bb
```