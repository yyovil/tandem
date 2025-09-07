{
  description = "flake to setup the dev env required for Tandem";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in
      {
        devShells.default =
          with pkgs;

          mkShell {
            buildInputs = [
              go
              gopls
              gotools
              vagrant
              sqlc
              goreleaser
            ];

            shellHook = ''
              export GOROOT="${pkgs.go}/share/go"
              export GOPATH="$PWD/.go"
              export PATH="$GOPATH/bin:$PATH"
              export GOBIN="$GOPATH/bin"
              source tandem.env
            '';
          };
      });
}