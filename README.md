#concordance-go

Implementation of concordance task for DataSift, implemented in go.

###Design

The design is based around tokenisers, take an reader (incoming set of bytes) and a channel of strings. The incoming bytes are parsed (as unicode), tokenised, and the tokens are passed out though the channel. This approach allows us to compose a set of tokenisers, producing the desired results.

By using channels we also allow ourselves to parallelise the algorithm (although currently it is essentially sequential). 
