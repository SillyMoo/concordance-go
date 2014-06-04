#concordance-go

Implementation of concordance task for DataSift, implemented in go.

###Design overview

The design is based around tokenisers, which take a reader (incoming set of bytes) and a channel of strings. The incoming bytes are parsed (as unicode), tokenised, and the tokens are passed out though the channel. Tokenisers are composable through a Compose method.

By using channels we also allow ourselves to parallelise the algorithm (although currently it is essentially sequential). 
