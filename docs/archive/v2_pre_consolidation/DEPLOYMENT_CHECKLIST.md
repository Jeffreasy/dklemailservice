# V1_47 Database Optimalisatie - Deployment Checklist

**Status**: Fix gepushed naar GitHub (commit b3d8286)  
**Datum**: 30 oktober 2025  
**Verwachte deployment tijd**: 3-5 minuten

---

## ✅ Deployment Status Tracking

### ✅ VOLTOOID
- [x] Database analyse uitgevoerd (33 tabellen)
- [x] V1_47 migratie gemaakt (30+ indexes)
- [x] Documentatie compleet (5 docs + 3 scripts)
- [x] Git commit & push (commit d0a348b)
- [x] Deployment fout gedetecteerd (order_number kolom)
- [x] Fix geïmplementeerd (Section 7 verwijderd)
- [x] Fix gepushed naar GitHub (commit b3d8286)

### 🔄 IN UITVOERING
- [ ] Render detecteert nieuwe push
- [ ] Render bouwt nieuwe versie
- [ ] Render deploy nieuwe versie
- [ ] V1_47 migratie wordt automatisch uitgevoerd

### ⏳ NOG TE DOEN (Na Deployment)
- [ ] Verificatie in Render logs
- [ ] ANALYZE commando uitvoeren
- [ ] Performance validatie

---

## 🔍 Monitor Render Deployment (NU)

### Stap 1: Check Render Dashboard

Ga naar: **https://dashboard.render.com/**

1. Klik op je DKL Email Service
2. Ga naar **"Events"** tab
3. Wacht tot je ziet:

```
✓ Deploying commit b3d8286...
✓ Building...
✓ Build successful
✓ Deploy live
```

**Verwachte tijd**: 3-5 minuten

---

### Stap 2: Monitor Logs voor V1_47 Migratie

Zodra deployment live is:

1. Ga naar **"Logs"** tab in Render Dashboard
2. **Zoek naar deze exacte regels:**

```json
{"niveau":"INFO","bericht":"Migratie wordt uitgevoerd","file":"V1_47__performance_optimizations.sql"}
{"niveau":"INFO","bericht":"Migratie succesvol uitgevoerd","file":"V1_47__performance_optimizations.sql"}
```

✅ **ALS JE DIT ZIET**: Migratie is SUCCESVOL! Ga naar Stap 3.

❌ **ALS JE ERROR ZIET**: Kopieer de error en laat het me weten.

---

### Stap 3: Update Database Statistieken (KRITIEK!)

**Na succesvolle migratie MOET je dit doen:**

```bash
# Option 1: Via PowerShell (lokaal)
$env:DATABASE_URL="postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"

psql "$env:DATABASE_URL" -c "ANALYZE;"

# Option 2: Via bash/WSL
export DATABASE_URL="postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"

psql "$DATABASE_URL" -c "ANALYZE;"
```

**Verwachte output:**
```
ANALYZE
```

Dit update de query planner statistieken zodat PostgreSQL de nieuwe indexes optimaal gebruikt!

---

### Stap 4: Verificatie van Indexes (Optioneel)

```bash
# Check hoeveel indexes er nu zijn
psql "$DATABASE_URL" -c "
SELECT 
    tablename, 
    COUNT(*) as index_count 
FROM pg_indexes 
WHERE schemaname = 'public' 
GROUP BY tablename 
HAVING COUNT(*) > 2
ORDER BY index_count DESC;"
```

**Je zou moeten zien:**
- `verzonden_emails`: 8-10 indexes (was: 2-3)
- `chat_messages`: 6-8 indexes (was: 2-3)
- `contact_formulieren`: 5-7 indexes (was: 1-2)
- `aanmeldingen`: 5-7 indexes (was: 1-2)

---

## 🎯 Performance Verificatie

### Test 1: Dashboard Query (Contact Forms)

```bash
psql "$DATABASE_URL" -c "
EXPLAIN ANALYZE
SELECT * FROM contact_formulieren 
WHERE status = 'nieuw' AND beantwoord = FALSE 
ORDER BY created_at DESC 
LIMIT 20;"
```

**Zoek naar**: `Index Scan` (GOED) vs `Seq Scan` (SLECHT)  
**Verwachte tijd**: < 10ms (was: ~50ms)

### Test 2: Email Tracking Query

```bash
psql "$DATABASE_URL" -c "
EXPLAIN ANALYZE
SELECT ve.*, cf.naam 
FROM verzonden_emails ve
LEFT JOIN contact_formulieren cf ON ve.contact_id = cf.id
WHERE ve.verzonden_op > NOW() - INTERVAL '7 days'
LIMIT 100;"
```

**Verwachte tijd**: < 20ms (was: ~100ms)

### Test 3: Full-Text Search

```bash
psql "$DATABASE_URL" -c "
EXPLAIN ANALYZE
SELECT * FROM contact_formulieren
WHERE to_tsvector('dutch', naam || ' ' || email || ' ' || bericht) @@ to_tsquery('dutch', 'test')
LIMIT 20;"
```

**Verwachte tijd**: < 30ms (was: ~500ms)

---

## 📊 Verwachte Resultaten

### Indexes Aangemaakt (V1_47)

| Categorie | Aantal | Impact |
|-----------|--------|--------|
| **FK Indexes** | 7 | 🔴 HIGH (50-90% sneller JOINs) |
| **Compound Indexes** | 5 | 🟡 MEDIUM (60-80% sneller dashboards) |
| **Partial Indexes** | 8 | 🟡 MEDIUM (70-85% sneller filters) |
| **FTS Indexes** | 3 | 🟢 HIGH (90-95% sneller search) |
| **Single Column** | 10 | 🟢 LOW-MEDIUM |

**Totaal**: ~33 nieuwe indexes

### Performance Verbetering

| Query Type | Voor | Na | Verbetering |
|------------|------|-----|-------------|
| JOIN queries | ~100ms | ~10ms | **90% sneller** |
| Dashboard | ~50ms | ~8ms | **84% sneller** |
| Search | ~500ms | ~25ms | **95% sneller** |
| Filters | ~60ms | ~12ms | **80% sneller** |

---

## ⚠️ Troubleshooting

### Als deployment faalt opnieuw:

1. **Check Render logs** voor de exacte foutmelding
2. **Kopieer de error** en laat het me weten
3. Ik kan dan direct de juiste fix maken

### Als ANALYZE faalt:

```bash
# Probeer per tabel
psql "$DATABASE_URL" -c "ANALYZE gebruikers;"
psql "$DATABASE_URL" -c "ANALYZE contact_formulieren;"
psql "$DATABASE_URL" -c "ANALYZE aanmeldingen;"
psql "$DATABASE_URL" -c "ANALYZE verzonden_emails;"
```

### Als queries nog steeds traag zijn:

```bash
# Check of indexes daadwerkelijk worden gebruikt
psql "$DATABASE_URL" -c "
SELECT * FROM pg_stat_user_indexes 
WHERE schemaname = 'public' 
  AND indexrelname LIKE 'idx_%'
ORDER BY idx_scan DESC;"
```

Als `idx_scan = 0` voor nieuwe indexes: Run ANALYZE opnieuw!

---

## 🎉 Success Criteria

Je deployment is succesvol als:

✅ Render logs tonen: `"Migratie succesvol uitgevoerd","file":"V1_47__performance_optimizations.sql"`  
✅ Applicatie start zonder errors  
✅ ANALYZE commando succesvol uitgevoerd  
✅ Dashboard queries gebruiken Index Scans (via EXPLAIN)  
✅ Geen performance degradation in Render metrics

---

## 📞 Next Steps Na Succesvolle Deployment

1. ✅ Monitor Render Metrics voor 24 uur
2. ✅ Check voor slow queries via pg_stat_statements (zie DATABASE_QUICK_REFERENCE.md)
3. ✅ Setup wekelijkse VACUUM (zie database/scripts/vacuum_analyze.sql)
4. ✅ Plan maandelijkse data cleanup (zie database/scripts/data_cleanup.sql)
5. ⏳ Overweeg partitioning als tabellen > 1M rows (zie database/scripts/setup_partitioning.sql)

---

**Deployment Tracking**:  
- Push tijd: 23:51 UTC  
- Commit: b3d8286  
- Status: Wachtend op Render deployment...

**Check status over 5 minuten** in Render Dashboard!