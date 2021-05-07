# PEGman

_An experimental, opinionated parser generator_


Maintaining the handcrafted parser for our [ntt project](https://github.com/nokia/ntt/)
is challenging, mostly because aligning parser-dependent syntax trees and
visitors is tedious. Hence I decided to automate this task with a custom
parser generator.  
Existing solutions did not fit our requirements very well. Some parsers are not
sufficient for large scale input. Others could not produce a lossless syntax
tree.

This parser generator is an experiment. I do not know where this project will
take us, yet.
