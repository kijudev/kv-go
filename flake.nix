{
  description = "Golang nix dev flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          nixd
          nil
          package-version-server
          cloc

          go
          gopls
          gotools
          go-tools
        ];

        shellHook = ''
          echo "==== Golang DevShell ===="
        '';
      };
    };
}
