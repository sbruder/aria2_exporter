{ config, lib, pkgs, ... }:
let
  cfg = config.services.aria2_exporter;
in
{
  options.services.aria2_exporter = {
    enable = lib.mkEnableOption "the prometheus aria2 exporter";
    package = lib.mkOption {
      type = lib.types.package;
      default = import ./default.nix { inherit pkgs; };
      example = "pkgs.aria2_exporter-fork";
      description = "The package to use for aria2_exporter";
    };
    rpcUrl = lib.mkOption {
      type = lib.types.str;
      default = "http://localhost:6800";
      example = "http://aria2.example.com:6800";
      description = "The RPC endpoint of aria2 aria2_exporter should connect to.";
    };
    rpcSecret = lib.mkOption {
      type = lib.types.str;
      default = "";
      example = "totallysecretsecret";
      description = "The RPC secret of aria2.";
    };
    listenAddress = lib.mkOption {
      type = lib.types.str;
      default = ":9578";
      example = "localhost:9578";
      description = "The address aria2_exporter should listen on.";
    };
  };

  config = {
    systemd.services.aria2_exporter = lib.mkIf cfg.enable {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      environment = {
        ARIA2_EXPORTER_LISTEN_ADDRESS = cfg.listenAddress;
        ARIA2_URL = cfg.rpcUrl;
        ARIA2_RPC_SECRET = cfg.rpcSecret;
      };
      serviceConfig = {
        ExecStart = "${cfg.package}/bin/aria2_exporter";
        Restart = "always";

        # systemd-analyze --no-pager security aria2_exporter.service
        CapabilityBoundingSet = null;
        DynamicUser = true;
        PrivateDevices = true;
        PrivateUsers = true;
        ProtectHome = true;
        RestrictAddressFamilies = [ "AF_INET" "AF_INET6" ];
        RestrictNamespaces = true;
        SystemCallArchitectures = "native";
        SystemCallFilter = "@system-service";
      };
    };
  };
}
