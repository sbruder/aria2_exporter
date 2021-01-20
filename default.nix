{ pkgs ? import <nixpkgs> { } }:
let
  gitignoreSrc = pkgs.fetchFromGitHub {
    owner = "hercules-ci";
    repo = "gitignore";
    rev = "c4662e662462e7bf3c2a968483478a665d00e717";
    sha256 = "1npnx0h6bd0d7ql93ka7azhj40zgjp815fw2r6smg8ch9p7mzdlx";
  };
  inherit (import gitignoreSrc { inherit (pkgs) lib; }) gitignoreSource;
in
pkgs.buildGoModule rec {
  name = "aria2_exporter";

  src = gitignoreSource ./.;

  subPackages = [ "." ];

  vendorSha256 = "1m0s6kp2pj3g54vyb7jzs1whcwyycs7v2p1malm3b8hw2jjp25xf";

  doCheck = false; # no tests

  meta.license = pkgs.lib.licenses.mit;
}
