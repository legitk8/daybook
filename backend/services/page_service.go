package services

import "daybook/backend/models"

var pages []models.Page

var nextID int = 1

func CreatePage(title string, content string) models.Page {
	page := models.Page{
		ID:      nextID,
		Title:   title,
		Content: content,
	}

	nextID++
	pages = append(pages, page)

	return page
}

func GetPageByID(id int) (models.Page, bool) {
	for _, page := range pages {
		if page.ID == id {
			return page, true
		}
	}
	return models.Page{}, false
}

func ListPages() []models.Page {
	return pages
}

func DeletePage(id int) bool {
	for i, page := range pages {
		if page.ID == id {
			pages = append(pages[:i], pages[i+1:]...)
			return true
		}
	}
	return false
}
