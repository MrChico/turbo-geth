let
  nixpkgs = import (builtins.fetchGit {
    url = "https://github.com/nixos/nixpkgs";
    ref = "release-20.03";
    rev = "5272327b81ed355bbed5659b8d303cf2979b6953";
  }) {};
in with nixpkgs;
buildGoModule rec {
  pname = "turbo-geth";
  version = "2020.08.04-alpha";

  src = fetchFromGitHub {
    owner = "ledgerwatch";
    repo = pname;
    rev = version;
    sha256 = "0mandbw06h4pigdk91vnw402agq5ynsjxmgkmnx3s9dqr6fhaimd";
  };

  modSha256 = "1fpnzhhy16kn27lvl8bbif080q05ql71znmdlq9ffddjlb8rwaad";


subPackages = [
"cmd/geth"
];

  meta = with stdenv.lib; {
    homepage = "https://geth.ethereum.org/";
    description = "Official golang implementation of the Ethereum protocol";
    license = with licenses; [ lgpl3 gpl3 ];
    maintainers = with maintainers; [ adisbladis lionello xrelkd ];
  };
}
