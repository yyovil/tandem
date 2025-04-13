# nixos guys ain't gotta do nothing. yall good. rest of them are just on their own for now.
{
  description = "flake to setup the dev env required for Tandem";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default =
        with pkgs;
        mkShell {
          buildInputs = [
            python311
            uv
            go
            # docker_25
            # docker is going to need a little config to work as expected. need to include in some specific group
          ];
          shellHook = ''
          # include dockerd and docker start cmds.
          # include cmds to install the deps for both the projects.
          source agents/.venv/bin/activate
          source .env

          '';
        };
    };
}
