{
  self,
  lib,
  buildGoApplication,
  go,
  installShellFiles,
}:

buildGoApplication {
  pname = "aio";
  version = import ./version.nix;
  src = self;
  pwd = self;
  modules = ./gomod2nix.toml;

  doCheck = false;

  ldflags = [
    "-s"
    "-w"
  ];

  inherit go;

  # Prevent Go from automatically downloading newer toolchains
  preBuild = ''
    export GOTOOLCHAIN=local
  '';

  nativeBuildInputs = [ installShellFiles ];
  postInstall = ''
    installShellCompletion --zsh ./completions/_aio
    installShellCompletion --bash ./completions/aio.bash
    installShellCompletion --fish ./completions/aio.fish
  '';

  meta = with lib; {
    description = "Aio is an open-source framework for augmenting humans using AI. It provides a modular framework for solving specific problems using a crowdsourced set of AI prompts that can be used anywhere";
    homepage = "https://github.com/creativerezz/aio";
    license = licenses.mit;
    platforms = platforms.all;
    mainProgram = "aio";
  };
}
