package v1

//func TestTaskRoutes_CreateTask(t *testing.T) {
//	testCases := []struct {
//		name string
//		url  string
//	}{
//		{},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			taskService := mock_service.NewMockTask(ctrl)
//			log := mock_logger.NewMockLogger(ctrl)
//
//			handler := NewHandler(services, nil, "")
//
//			r := gin.Default()
//			r.POST(url+"/users/report", handler.CreateCSVReportAndURL)
//
//			requestBody := map[string]interface{}{
//				"date": expectedDate,
//			}
//
//			jsonBody, err := json.Marshal(requestBody)
//			require.NoError(t, err)
//
//			w := httptest.NewRecorder()
//
//			ctx := context.Background()
//			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/report", bytes.NewBuffer(jsonBody))
//			require.NoError(t, err)
//			req.Header.Set("Content-Type", "application/json")
//
//			r.ServeHTTP(w, req)
//
//			require.Equal(t, http.StatusOK, w.Code)
//
//			var responseBody map[string]interface{}
//			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
//			require.NoError(t, err)
//
//			actualURL, ok := responseBody["report_url"]
//			require.Equal(t, expectedUrl, actualURL)
//			require.True(t, ok)
//		})
//	}
//}
