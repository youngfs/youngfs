package fullpath

type PathLink struct {
	Name string
	Link string
}

func (fp FullPath) ToPathLink() []PathLink {
	list := fp.Split()

	pathLinks := make([]PathLink, 0)
	link := ""
	for _, str := range list {
		link += string(str) + "/"
		u := PathLink{
			Name: string(str) + "/",
			Link: link,
		}
		pathLinks = append(pathLinks, u)
	}
	return pathLinks
}
