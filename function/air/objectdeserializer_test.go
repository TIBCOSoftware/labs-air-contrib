package air

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFNObjectDeserializer_Eval(t *testing.T) {
	f := &fnObjectDeserializer{}
	objectstr := "o2phcGlWZXJzaW9uYnYyaXJlcXVlc3RJZHgkMTIwMmVjN2ItYzViYS00NWY4LThhZTItYTQwMTExZDkyYzU1ZWV2ZW50p2phcGlWZXJzaW9uYnYyYmlkeCRiNGQ0YjFlMS1hNmY5LTRmODUtYTUxNi1iM2JmNWEyYjVmZjhqZGV2aWNlTmFtZWpSRVNURGV2aWNla3Byb2ZpbGVOYW1lc0dlbmVyaWMtUkVTVC1EZXZpY2Vqc291cmNlTmFtZW1pbWFnZV9yZWFkaW5nZm9yaWdpbhsW5T8QCkijtGhyZWFkaW5nc4GoYmlkeCQxNjJiOWM0MS02NTllLTRlYWUtODU2NC1mODBjZmU5ZDA5YTdmb3JpZ2luGwAAAYAfZVZxamRldmljZU5hbWVqUkVTVERldmljZWxyZXNvdXJjZU5hbWVtaW1hZ2VfcmVhZGluZ2twcm9maWxlTmFtZXNHZW5lcmljLVJFU1QtRGV2aWNlaXZhbHVlVHlwZWZCaW5hcnlrYmluYXJ5VmFsdWVZBfP/2P/gABBKRklGAAEBAQBIAEgAAP/bAEMA///////////////////////////////////////////////////////////////////////////////////////bAEMB///////////////////////////////////////////////////////////////////////////////////////AABEIAOoBOQMBIgACEQEDEQH/xAAXAAEBAQEAAAAAAAAAAAAAAAAAAQID/8QAJBABAQEAAgEEAwEBAQEAAAAAAAERITFBAhJRcWGBkbGhwfD/xAAVAQEBAAAAAAAAAAAAAAAAAAAAAf/EABYRAQEBAAAAAAAAAAAAAAAAAAARAf/aAAwDAQACEQMRAD8AyuIvIIqCDWBDAKCzgVMw52616mNBqT4XMxJWqDFl1F088gi8flbPMZ0GuKmG1N4BeDZmMoDWmooNQZlxrd+gF6TyAW/wh7d8p0DWpvyzuGgvp4tW3nUWwCeqXdS3UzsBZt+UswOQQxf8Xr9iMr+kAMVABQ0DpqVhqAqTvab+P2ngGrZ8MNJgos1c80ugcT81PpFgLfVwnFnfJZUz5oFiL9LZmCICydgnHkPbW8+aqubc9N+lBDDIALwlm9IoOeLjVmsy5RVWX5TP3D2//VBatmxnmcU8dgf4lnK6t5BIXkz58nQIB/BEBVEFADUMBZVRrKglRagOmzGe0+y54FXj5XGF3IC2/wCM4sakgMwtt4a9sXiKjMl+v9a6QBUAEBcBBcQAABLNUBl07ZY2it2b56YWXtdxBI1mIW2/QL3D3SeEl/iYDW74ZwkXAjICov8A6i/CAY19JGp0is1rb0zSAoICiNT08AyNZelyQQkU0UQABcRQEAFjPqtnS6t5Bz2/K31fEMXZABe0AFxACzQBnLCS+W4lvKKmJ9qznYNzLGOVnB4A7TGuL+Fz7Bk4Qoi7qdLJ5OhVnJZfC8KDneBbzVBDtpM0Cfl0Zkk+zVRdRFABQQwJdAtkT3c9FnmE9OA0ytsjMoKbnYtnAKz7YzzDaDeyM+5kBdvgiwwBQALNVQYzj8r6Z3q8RdRTIzk8NHQjFh7qtqaKHRU0GtLjGio3xDUEVCTfKrICyKiKgAAqANIkvj+NAifmKAus+q065/qg5rFsygCW1QCGEigmKKCCs+74Bo4jG0BfcltoAhovCLi31fhPd+DYAW2pytRRqs/bTNQMCaqgTRZEFk81VRUQAAABUALyS+KHYNIrE50GqztnSXSAjSYoALgC4cTtm0GuIz7vjhAAAAAAEAWy94i20EgLJqKS/LX6v8J6fmtZAcgFBrisgNGpKXmoNarOKqAAAABMDsGkZ2pdvYLb4n9SAC1FAFJC2QFZvq+EttAAAAAAQFEAAAABRub4Zk10/EA/6qXZOPDHv/AiBnyIq8GT4ImgJeKL2Bqpn+qqKEAAARUAEaQEVTAF+02RndBbfhAAAAEUBABUAAAAABFWSXyCz/jeyTWMvjku8cCt6mT4T1dMbfmguqyYgqYsATKAqDTLUAVAFRfCAAACgELfhLdQAAAEBUAAAAAAEBRABUAF1FBrU9xjKK3bPnwzwLii4zY1ek4QZ5a01lRq2IAigAoigKgCgbgL0xbp2AigAgAAAAAAAIpYCAvQIAAAAuIA0IqKg0uAzUKsuKGZNp+S3UQXdJnllVRbiacgq8CLl+ERRFUNZVAFEAAAAAAABAVFXIDK6WACoIqUXwiouwSAL2YagL4i6yoKu/lkRRU6TVG4l/BIs46QYWLhc8Khhis2Ipq6iwEaL0gKy0zVQAAAAAAVDRVxMEBdQaAioagqEoDNFFRBQDDGhKsZGmRAD+qKiw4v2gmm0xFFVlrRVSpKIKEToGoyALCoqoiiAAgCooIqAqiiCDSUEW4gqAAAAo1GV/aCqx+2wS1lagACoAAIqgyoAB5Kiqnf2QvYCos6UTysQEVABFgAuGKCs0L2AaaigunhPB4RAUUEUAQKAigC7SJ8gCoAAD//2WltZWRpYVR5cGVqaW1hZ2UvanBlZw=="
	v, err := function.Eval(f, []byte(objectstr), nil)
	assert.Nil(t, err)
	fmt.Println("#### ", v)
}

func TestFNObjectDeserializer_Eval2(t *testing.T) {
	f := &fnObjectDeserializer{}
	objectstr := "{\"sample_type\":\"JSON\"}"
	v, err := function.Eval(f, []byte(objectstr), nil)
	assert.Nil(t, err)
	fmt.Println("#### ", v)
}
