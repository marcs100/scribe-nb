package ui

type PageView struct{
	NumberOfPages int
	CurrentPage int
}

func (pv *PageView) PageForward() int{
	if pv.CurrentPage < pv.NumberOfPages{
		pv.CurrentPage += 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageView) PageBack() int{
	if pv.CurrentPage > 1{
		pv.CurrentPage -= 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageView) Reset(){
	pv.CurrentPage = 0
	pv.NumberOfPages = 0
}

