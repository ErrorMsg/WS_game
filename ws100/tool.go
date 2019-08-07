//tools

package srv



func SliceDel(s interface{}, n int) interface{
	r := s[n]
	s = s[:n]+s[n,len(s)-1]
	return r
}