# nixos guys ain't gotta do nothing. yall good. rest of them are just on their own for now.
{
  description = "flake to setup the dev env required for Tandem";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; config.allowUnfree = true; };
    in
    {
      devShells.${system}.default =
        with pkgs;
        mkShell {
          buildInputs = [
            python311
            ruff
            uv
            go
            gopls
            gotools
            vagrant
            starship

            # TODO:
            # better dev setup for go
            # docker_25
            # docker is going to need a little config to work as expected. need to include in some specific group
          ];
          shellHook = ''
            export STARSHIP_CONFIG=$PWD/.config/starship.toml
            export GOROOT="${pkgs.go}/share/go"
            export GOPATH="$PWD/.go"
            export PATH="$GOPATH/bin:$PATH"
            export GOBIN="$GOPATH/bin"
            source agents/.venv/bin/activate
            source monorepo.sh
            source tui/.env
            eval "$(starship init bash)"

            # TODO:
            # include dockerd and docker start cmds.
            # include cmds to install the deps for both the projects.
          '';
        };
    };
}
