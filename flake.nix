{
  description = "Bluetooth Advertisement Sniffer - A tool for capturing and analyzing BLE advertisements";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        bluetooth-sniffer = pkgs.buildGoModule {
          pname = "bluetooth-sniffer";
          version = "1.0.0";

          src = ./.;
          vendorHash = "sha256-CjEFbHnf6x1ebrIuWiNelEWI7E2CC4vrQo86gJ1+wQU=";

          meta = with pkgs.lib; {
            description = "Bluetooth Low Energy advertisement data sniffer";
            homepage = "https://github.com/nxm/bluetooth-sniffer";
            license = licenses.mit;
            platforms = platforms.linux ++ platforms.darwin;
          };

          # BlueZ is the Linux BLE stack; on Darwin the cbgo backend links
          # against the system CoreBluetooth framework instead.
          buildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.bluez ];
        };
      in
      {
        packages = {
          default = bluetooth-sniffer;
          bluetooth-sniffer = bluetooth-sniffer;
        };

        apps.default = flake-utils.lib.mkApp {
          drv = bluetooth-sniffer;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
            golangci-lint
          ] ++ lib.optionals stdenv.isLinux [ bluez ];
        };
      }
    );
}