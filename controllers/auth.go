package controllers

func LoginOrRegister(c *gin.Context) {  // Get model if exist
	var user models.Client
  
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	  return