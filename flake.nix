{
  description = "An aria2 Exporter for Prometheus";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs }: {
    overlay = final: prev: {
      aria2_exporter = prev.buildGoModule rec {
        name = "aria2_exporter";

        src = self;

        subPackages = [ "." ];

        vendorSha256 = "1m0s6kp2pj3g54vyb7jzs1whcwyycs7v2p1malm3b8hw2jjp25xf";

        doCheck = false; # no tests

        meta = with prev.lib; {
          license = licenses.mit;
          maintainer = with mainatiners; [ sbruder ];
        };
      };
    };

    nixosModules.aria2_exporter = {
      imports = [ ./module.nix ];

      nixpkgs.overlays = [
        self.overlay
      ];
    };
  } // flake-utils.lib.eachSystem [ "aarch64-linux" "x86_64-linux" ] (system:
    let
      pkgs = import nixpkgs { inherit system; overlays = [ self.overlay ]; };
    in
    rec {
      packages = {
        inherit (pkgs) aria2_exporter;
      };

      defaultPackage = packages.aria2_exporter;

      checks = {
        integration-test = import ./test.nix {
          inherit nixpkgs system;
          inherit (self) nixosModules;
        };
      };
    });
}
