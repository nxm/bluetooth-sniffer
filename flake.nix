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
          vendorHash = "sha256-UB/V4CO1YxqjNvROhCJAyWa3q79YrvkWa1R2sdSf8Zo=";

          meta = with pkgs.lib; {
            description = "Bluetooth Low Energy advertisement data sniffer";
            homepage = "https://github.com/nxm/bluetooth-sniffer";
            license = licenses.mit;
            platforms = platforms.linux;
          };

          buildInputs = with pkgs; [
            bluez
          ];
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
            bluez
          ];
        };
      }
    );
}