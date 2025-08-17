{
  description = "flake to setup the dev env required for Tandem";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        config.allowUnfree = true;
      };
    in
    {
      devShells.${system}.default =
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
          '';
        };
    };
}