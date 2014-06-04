#concordance-go

Implementation of concordance task for DataSift, implemented in go.

###Design overview

The design is based around tokenisers, which take a reader (incoming set of bytes) and a channel of strings. The incoming bytes are parsed (as unicode), tokenised, and the tokens are passed out though the channel. Tokenisers are composable through a Compose method.

By using channels we also allow ourselves to parallelise the algorithm (although currently it is essentially sequential). 

###Building

Assuming you have go set up, and a ```GOPATH``` variable defined, check this code out to ```$GOPATH/src/gihub.com/SillyMoo/concordance-go```.  
- To build type ```go install github.com/SillyMoo/concordance-go```
- To test type ```go test github.com/SillyMoo/concordance-go github.com/SillyMoo/concordance-go/tokeniser```
- To run use $GOPATH/bin/concordance-go and pipe in some text to be analysed


Options:  
- --outputFormat=[standard|wordIdx] - defines the output format, standard (the default) will output in the standard format as per the assignment pdf. wordIdx will output in the format ```token {frequency:sp1.wp1,sp2.wp2,...}``` where sp is the index of the sentence in which the word occurs, and wp is the index of the word within sentence.  
