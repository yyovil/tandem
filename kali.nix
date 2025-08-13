{ }:

let
  system = "x86_64-linux";
  pkgs = import <nixpkgs> { inherit system; };
  kali-base = pkgs.dockerTools.pullImage {
    imageName = "kalilinux/kali-last-release";
    imageDigest = "sha256:d44a3c0423addaaae40a04f8c935e245067688591cd78545504ea70802ed40ba";
    sha256 = "sha256-CYe1AMOZlEuvZfrOmpfHjjdmxypon7QXC3lG/ibACHI=";
    finalImageName = "kali";
    finalImageTag = "latest";
  };
in
pkgs.dockerTools.buildLayeredImage {
  fromImage = kali-base;
  name = "kali";
  tag = "withtools";
  contents = with pkgs; [
    nettools
    iproute2
    nmap
    dirb
    metasploit
    nano
    nikto
    exploitdb
    enum4linux
    python3
    sqlmap
  ];
  config = {
    Cmd = [ "bash" ];
    WorkingDir = "/home/tandem";
  };
}
