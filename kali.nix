{ }:
let
  system = "x86_64-linux";
  pkgs = import <nixpkgs> { inherit system; };
in
pkgs.dockerTools.buildLayeredImage {
  fromImage = "kali:latest";
  name = "kali";
  tag = "withtools";
  contents = with pkgs; [
    nmap
    dirb
    metasploit
  ];
}
