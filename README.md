#concordance-go

Implementation of concordance task for DataSift, implemented in go.

###Design overview

This is a much faster version of concordance-go. It's still a work in progres and may well always be. Why? Well it's ugly code, I've taken all the good design elements of the original and instead gone for one very fast monolithic function that does all the tokenisation in one fell swoop. It's fast, but boy readable it sure ain't.

So how fast? Well this was the last timing I had for the normal verison on the master branch (running through 5.5mb worth of the complete works of shakespear):

real    0m1.888s  
user    0m1.841s  
sys     0m0.090s  

And here is the new version:

real    0m0.805s  
user    0m0.764s  
sys     0m0.050s  

Yeah, that fast, it processes roughly 6.6mb a second, not too shabby.

It has issues though, beyond just the ugliness, it currently does not handle carriage returns correctly, so it's not handling the bard quite right yet.
