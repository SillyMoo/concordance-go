#concordance-go

Implementation of concordance task for DataSift, implemented in go.

###Design overview

The design is based around tokenisers, which take a reader (incoming set of bytes) and a channel of strings. The incoming bytes are parsed (as unicode), tokenised, and the tokens are passed out though the channel. Tokenisers are composable through a Compose method.

The use of channels is not strictly necessary, but does afford easy composition. The go routine and channel composition makes producer/consumer type abstractions very simple, and a chain of tokenisers fits this pattern perfectly.

####Performance considerations
Channels do add a performance overhead, if performance was a primary consideration you could move to using the bufio scanner and splitfunc (as used inside the tokenisers) directly. You could also avoid using a change of tokenisers and do splitting lines and tokens in a single function (which would reduce the GC impact described below).

Note that since we are only dependent upon a single input and output resource (standard input and standard output) the design has limited concurrency (i.e thread blocking due to resource contention is not an issue). A more concurrent approach would be possible as the tokeniser is based around the channel abstraction, but the additional complexity was not justified by the current requirements.

There is no real effort to limit the amount of strings or byte arrays created, this could be an issue if processing very large files (causing GC pauses), but that does not seem to be the intent of the assignment so I left it as is for now. Also to order the results a slice containing all tokens is created, again for a very large input this could be problematic.

###Building
First if you have not already done so install [go](http://golang.org/doc/install). You will need to set an environmental variable called GOPATH after installing go, this a path to the 'workspace' in which go code is developed and installed, see [here](http://golang.org/doc/code.html#GOPATH) for more details
Once go is setup you can get the code by type ```go get github.com/SillyMoo/concordance-go```.
- To build type ```go install github.com/SillyMoo/concordance-go```
- To test type ```go test github.com/SillyMoo/concordance-go github.com/SillyMoo/concordance-go/tokeniser```
- To run use $GOPATH/bin/concordance-go and pipe in some text to be analysed


Options:  
- --outputFormat=[standard|wordIdx] - defines the output format, standard (the default) will output in the standard format as per the assignment pdf. wordIdx will output in the format ```token {frequency:sp1.wp1,sp2.wp2,...}``` where sp is the index of the sentence in which the word occurs, and wp is the index of the word within sentence.  


###Addtional performance notes
As a quick check to confirm that the additional concurrency would not improve performance I created a concordance from the complete works of [shakespear](http://www.gutenberg.org/ebooks/100.txt.utf-8), piping the output to /dev/null. The output from 'time' is shown below:

real    0m1.953s  
user    0m1.902s  
sys     0m0.086s  

We can see that real closely matches user (~2.6% delta), showing that we are quite effectively using the single thread without too much blocking (some blocking is inevitable, we are reading off disk after all). With some code changes and raising GOMAXPROCS we could make a parallel implementation. But it is probably better to split the file up in a separate process, and then stand up multiple instances of concordance-go to process the file parts, finally merging the multiple concordances in a final step (thus we could, using this mechanism, scale over multiple VMs/servers).

Update:
Added a small buffer to all channels, as this showed a small performance improvement. Previously all channels were unbuffered which meant we could always process a word whilst waiting for a word to be pulled of the standard input (i.e if the go routines had finished processing all previously retrieved bytes).

This lead to a small but consistent performance improvement:

real    0m1.888s  
user    0m1.841s  
sys     0m0.090s  

