package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tinygo.org/x/bluetooth"
)

type BluetoothSniffer struct {
	adapter    *bluetooth.Adapter
	filterAddr string
}

func NewBluetoothSniffer(filterAddr string) (*BluetoothSniffer, error) {
	err := bluetooth.DefaultAdapter.Enable()
	if err != nil {
		return nil, fmt.Errorf("failed to enable bluetooth: %w", err)
	}

	return &BluetoothSniffer{
		adapter:    bluetooth.DefaultAdapter,
		filterAddr: strings.ToLower(filterAddr),
	}, nil
}

func (bs *BluetoothSniffer) Start(ctx context.Context) error {
	fmt.Printf("Starting Bluetooth advertisement scanner...\n")
	if bs.filterAddr != "" {
		fmt.Printf("Filtering devices containing address: %s\n", bs.filterAddr)
	}
	fmt.Printf("Press Ctrl+C to stop\n\n")

	scanDone := make(chan error, 1)

	go func() {
		err := bs.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			deviceAddr := strings.ToLower(result.Address.String())

			if bs.filterAddr != "" && !strings.Contains(deviceAddr, bs.filterAddr) {
				return
			}

			bs.displayAdvertisementData(result)
		})
		scanDone <- err
	}()

	select {
	case err := <-scanDone:
		if err != nil {
			return fmt.Errorf("scanning failed: %w", err)
		}
	case <-ctx.Done():
		err := bs.adapter.StopScan()
		if err != nil {
			return fmt.Errorf("failed to stop scanning: %w", err)
		}
	}

	return nil
}

func (bs *BluetoothSniffer) displayAdvertisementData(result bluetooth.ScanResult) {
	timestamp := time.Now().Format("15:04:05.000")

	fmt.Printf("=== [%s] ===\n", timestamp)
	fmt.Printf("Address: %s\n", result.Address.String())
	fmt.Printf("RSSI: %d dBm\n", result.RSSI)
	distance := rssiToDistance(int(result.RSSI))
	fmt.Printf("Estimated Distance: %.2fm\n", distance)

	localName := result.LocalName()
	if localName != "" {
		fmt.Printf("Local Name: %s\n", localName)
	}

	manufacturerData := result.ManufacturerData()
	if len(manufacturerData) > 0 {
		fmt.Printf("Manufacturer Data:\n")
		for _, data := range manufacturerData {
			fmt.Printf("  Company ID: 0x%04x\n", data.CompanyID)
			fmt.Printf("  Data: %s\n", hex.EncodeToString(data.Data))
		}
	}

	serviceData := result.ServiceData()
	if len(serviceData) > 0 {
		fmt.Printf("Service Data:\n")
		for _, data := range serviceData {
			fmt.Printf("  Service UUID: %s\n", data.UUID.String())
			fmt.Printf("  Data: %s\n", hex.EncodeToString(data.Data))
		}
	}

	rawData := result.Bytes()
	if len(rawData) > 0 {
		fmt.Printf("Raw Advertisement Data: %s\n", hex.EncodeToString(rawData))
	}

	fmt.Printf("\n")
}

// rssiToDistance converts RSSI to estimated distance
// Simple formula: distance = 10^((Tx Power - RSSI) / (10 * n))
// Where Tx Power ≈ -59 dBm at 1m, n ≈ 2 (path loss exponent)
func rssiToDistance(rssi int) float64 {
	if rssi == 0 {
		return 0
	}
	txPower := -59.0 // dBm at 1 meter
	pathLoss := 2.0
	distance := math.Pow(10, (txPower-float64(rssi))/(10*pathLoss))
	return math.Max(0.1, distance) // Minimum 0.1m
}

func main() {
	var filterAddr = flag.String("addr", "", "Filter devices by address (partial match, case insensitive)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Printf("\nShutting down...\n")
		cancel()
	}()

	sniffer, err := NewBluetoothSniffer(*filterAddr)
	if err != nil {
		log.Fatalf("Failed to create bluetooth sniffer: %v", err)
	}

	if err := sniffer.Start(ctx); err != nil {
		log.Fatalf("Failed to start sniffer: %v", err)
	}

	fmt.Println("Bluetooth sniffer stopped.")
}
