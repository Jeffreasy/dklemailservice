# 🚀 WebSocket Deployment Checklist

## ✅ Pre-Deployment Checklist

### Code Review

- [ ] **Backend code review**
  - [ ] [`services/steps_hub.go`](../services/steps_hub.go) - StepsHub implementation
  - [ ] [`handlers/steps_websocket_handler.go`](../handlers/steps_websocket_handler.go) - Handler
  - [ ] [`services/steps_service.go`](../services/steps_service.go) - Service updates
  - [ ] Error handling adequate
  - [ ] Logging comprehensive

- [ ] **Tests passed**
  - [x] `go test ./tests/steps_hub_test.go -v` → 8/8 PASSED ✅
  - [ ] Integration tests run
  - [ ] Manual testing completed

- [ ] **Security review**
  - [ ] JWT validation implemented
  - [ ] Permission checks in place
  - [ ] Input validation
  - [ ] No sensitive data in logs

### Code Integration

- [ ] **Update main.go**
  ```go
  // Add these lines:
  stepsHub := services.NewStepsHub(stepsService, gamificationService)
  stepsService.SetStepsHub(stepsHub)
  go stepsHub.Run()
  
  stepsWsHandler := handlers.NewStepsWebSocketHandler(stepsHub, authService)
  stepsWsHandler.RegisterRoutes(app)
  ```

- [ ] **Build test**
  ```bash
  go build -o dklemailservice.exe .
  # Should compile without errors
  ```

- [ ] **Run locally**
  ```bash
  go run main.go
  # Check logs for: "StepsHub started successfully"
  ```

### Environment Configuration

- [ ] **Update .env** (optioneel)
  ```bash
  WEBSOCKET_ENABLED=true
  WEBSOCKET_MAX_CONNECTIONS=10000
  ```

- [ ] **CORS settings**
  ```bash
  CORS_ORIGINS=https://your-frontend.com,wss://your-frontend.com
  ```

---

## 🧪 Testing Checklist

### Local Testing

- [ ] **WebSocket connection**
  ```bash
  wscat -c "ws://localhost:8080/ws/steps?user_id=test"
  # Should connect successfully
  ```

- [ ] **Subscribe to channels**
  ```json
  {"type":"subscribe","channels":["total_updates"]}
  # Should receive confirmation
  ```

- [ ] **Trigger update**
  ```bash
  curl -X POST http://localhost:8080/api/steps \
    -H "Authorization: Bearer TOKEN" \
    -d '{"steps": 100}'
  # WebSocket should receive broadcast
  ```

- [ ] **Check stats**
  ```bash
  curl http://localhost:8080/api/ws/stats \
    -H "Authorization: Bearer ADMIN_TOKEN"
  # Should show 1 client
  ```

### Frontend Testing

- [ ] **React component**
  - [ ] Copy `useStepsWebSocket.ts` to project
  - [ ] Import in component
  - [ ] Test connection
  - [ ] Verify updates received

- [ ] **Browser console test**
  ```javascript
  const ws = new WebSocket('ws://localhost:8080/ws/steps?user_id=test');
  ws.onmessage = (e) => console.log(JSON.parse(e.data));
  // Should log ping/pong and updates
  ```

---

## 🌐 Staging Deployment

### Pre-Deploy

- [ ] **Git commit**
  ```bash
  git add services/steps_hub.go
  git add handlers/steps_websocket_handler.go
  git add services/steps_service.go
  git add tests/steps_hub_test.go
  git commit -m "feat: Add WebSocket support for real-time steps tracking"
  ```

- [ ] **Push to staging**
  ```bash
  git push staging main
  ```

### Verify Staging

- [ ] **Check logs**
  ```bash
  # Heroku example
  heroku logs --tail --app your-app-staging
  
  # Should see:
  # INFO StepsHub started successfully
  ```

- [ ] **Test WebSocket**
  ```bash
  wscat -c "wss://your-app-staging.herokuapp.com/ws/steps?user_id=test&token=TOKEN"
  ```

- [ ] **Monitor metrics**
  ```bash
  curl https://your-app-staging.herokuapp.com/api/ws/stats
  ```

- [ ] **Load test** (optioneel)
  ```bash
  # 100 concurrent connections
  ./scripts/websocket_load_test.sh 100
  ```

### Staging Success Criteria

- [ ] WebSocket connects successfully
- [ ] Messages broadcast correctly
- [ ] No memory leaks (monitor voor 24h)
- [ ] No errors in logs
- [ ] Performance acceptable (< 50ms latency)

---

## 🚀 Production Deployment

### Pre-Production

- [ ] **Staging test passed** (minimum 24h monitoring)
- [ ] **Team approval**
- [ ] **Rollback plan ready**
- [ ] **Monitoring configured**

### Deployment Strategy

**Option 1: Direct Deploy (Simple)**
```bash
git push production main
```

**Option 2: Canary Deploy (Recommended)**
```bash
# 1. Deploy to 10% of instances
# 2. Monitor for 2 hours
# 3. Scale to 50%
# 4. Monitor for 4 hours
# 5. Full rollout
```

### Post-Deploy Verification

- [ ] **Health check**
  ```bash
  curl https://your-app.com/api/health
  ```

- [ ] **WebSocket connectivity**
  ```bash
  wscat -c "wss://your-app.com/ws/steps?user_id=test&token=TOKEN"
  ```

- [ ] **Stats check**
  ```bash
  curl https://your-app.com/api/ws/stats
  ```

- [ ] **Monitor logs** for errors
- [ ] **Monitor metrics**
  - Connection count
  - Message latency
  - Error rate
  - CPU/Memory usage

---

## 📊 Monitoring (First 48 Hours)

### Metrics to Watch

| Metric | Check Frequency | Alert Threshold |
|--------|----------------|-----------------|
| Connection Count | Every 5 min | > 8,000 |
| Message Latency | Every 1 min | > 100ms (p95) |
| Error Rate | Every 1 min | > 1% |
| CPU Usage | Every 5 min | > 80% |
| Memory Usage | Every 5 min | > 90% |

### Log Monitoring

**Watch for**:
- [ ] "StepsHub started successfully" at startup
- [ ] No "WebSocket auth failed" errors
- [ ] Connection/disconnection patterns normal
- [ ] No goroutine leaks

### User Feedback

- [ ] Survey early adopters
- [ ] Monitor support tickets
- [ ] Check for bug reports
- [ ] Collect feature requests

---

## 🔄 Rollback Procedure

### If Issues Detected

1. **Immediate rollback**
   ```bash
   git revert HEAD
   git push production main
   ```

2. **Or revert to previous version**
   ```bash
   git reset --hard <previous-commit-hash>
   git push --force production main
   ```

3. **Disable WebSocket** (zonder code rollback)
   ```bash
   # Set environment variable
   WEBSOCKET_ENABLED=false
   
   # Restart app
   ```

### Post-Rollback

- [ ] Notify team
- [ ] Document issue
- [ ] Fix in staging
- [ ] Re-test
- [ ] Re-deploy

---

## 📈 Success Metrics

### Week 1

- [ ] No critical bugs
- [ ] < 1% error rate
- [ ] User feedback positive
- [ ] Performance within targets

### Month 1

- [ ] 100+ active users
- [ ] Real-time updates working correctly
- [ ] No scaling issues
- [ ] User retention improved

### Quarter 1

- [ ] Feature widely adopted
- [ ] Engagement metrics up
- [ ] Performance optimized
- [ ] Ready for mobile app

---

## 🎯 Post-Deployment Tasks

### Week 1

- [ ] Monitor closely
- [ ] Fix any bugs
- [ ] Optimize performance if needed
- [ ] Document learnings

### Week 2-4

- [ ] Collect user feedback
- [ ] Plan enhancements
- [ ] Performance tuning
- [ ] Scale testing

### Month 2+

- [ ] Redis pub/sub voor multi-instance (if needed)
- [ ] Message persistence (if needed)
- [ ] Advanced analytics
- [ ] Mobile app integration

---

## 📝 Documentation Updates

- [ ] Update main README.md met WebSocket features
- [ ] Update API docs met WebSocket endpoints
- [ ] Create video tutorial (optioneel)
- [ ] Update changelog

---

## ✅ Final Sign-Off

### Technical Lead

- [ ] Code reviewed and approved
- [ ] Architecture validated
- [ ] Security verified
- [ ] Performance acceptable

### Product Owner

- [ ] Features meet requirements
- [ ] User experience approved
- [ ] Documentation adequate

### DevOps

- [ ] Deployment plan approved
- [ ] Monitoring configured
- [ ] Rollback tested
- [ ] Infrastructure ready

---

## 🎊 Ready to Deploy!

Wanneer alle checkboxen ✅ zijn:

1. **Merge to main**
2. **Deploy to production**
3. **Monitor closely**
4. **Celebrate success!** 🎉

---

**Document**: Deployment Checklist  
**Version**: 1.0  
**Status**: ✅ Ready for Use  
**Next**: Start deploying!