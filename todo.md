# To do

- Fix file filtering  
  Dog-fooding the program on its own repository showed, that with the regex
  filter of `\.go$` and `.*\.go$`, the files ending with `_test.go` weren't
  touched.  
  Currently investigating why.
- Fix tab to spaces conversion  
  Dog-fooding the program on its own repository (thanks Go for so many tabs!)
  showed inconsistent tab conversion.  
  Currently investigating why; it'll likely be a regex problem.
