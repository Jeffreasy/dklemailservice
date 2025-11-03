# 📱 Gamification API - Mobile Implementation Guide

> **Version:** 1.0 | **Platforms:** Android, iOS, React Native, Flutter | **Last Updated:** 2025-11-02

Complete mobile implementatie gids voor de Gamification API van DKL Email Service.

---

## 📋 Table of Contents

1. [Mobile Optimizations](#-mobile-optimizations)
2. [Android Implementation](#-android-implementation-kotlin)
3. [iOS Implementation](#-ios-implementation-swift)
4. [React Native Implementation](#-react-native-implementation)
5. [Flutter Implementation](#-flutter-implementation)
6. [Offline Support](#-offline-support)
7. [Performance Tips](#-performance-tips)
8. [Troubleshooting](#-troubleshooting)

---

## 🚀 Mobile Optimizations

De Gamification API is geoptimaliseerd voor mobile gebruik met:

### ✅ Implemented Features

- **Paginatie**: Leaderboard met page/limit voor data efficiency
- **Query Parameters**: Flexible filtering voor targeted data
- **JSON Format**: Lightweight responses
- **CORS Enabled**: Cross-origin support
- **Bearer Auth**: Standard JWT authentication

### 📱 Mobile-Specific Considerations

**Bandwidth Optimization:**
- Gebruik `limit` parameter om data te beperken (default: 50, max: 100)
- Filter met `route`, `year`, `top_n` voor specifieke data
- Public endpoints (`/api/leaderboard`, `/api/badges`) vereisen geen auth

**Caching Strategy:**
- Cache leaderboard lokaal voor 5 minuten
- Cache badges indefinitely (invalidate on badge change)
- Refresh achievements na steps update

**Background Sync:**
- Sync steps in background
- Check en award badges periodiek
- Update leaderboard tijdens off-peak

---

## 📱 Android Implementation (Kotlin)

### 1. Data Models

```kotlin
// File: models/Gamification.kt

data class LeaderboardEntry(
    val id: String,
    val naam: String,
    val route: String,
    val steps: Int,
    val achievement_points: Int,
    val total_score: Int,
    val rank: Int,
    val badge_count: Int,
    val joined_at: String
)

data class LeaderboardResponse(
    val entries: List<LeaderboardEntry>,
    val total_entries: Int,
    val current_page: Int,
    val total_pages: Int,
    val limit: Int
)

data class Badge(
    val id: String,
    val name: String,
    val description: String,
    val icon_url: String?,
    val criteria: BadgeCriteria,
    val points: Int,
    val is_active: Boolean,
    val display_order: Int,
    val created_at: String,
    val updated_at: String,
    val earned_count: Int? = null,
    val last_earned_at: String? = null,
    val earned_by_current_user: Boolean? = null
)

data class BadgeCriteria(
    val min_steps: Int? = null,
    val min_days: Int? = null,
    val consecutive_days: Int? = null,
    val route: String? = null,
    val early_participant: Boolean? = null,
    val has_team: Boolean? = null
)

data class Achievement(
    val id: String,
    val participant_id: String,
    val earned_at: String,
    val badge: Badge
)

data class ParticipantSummary(
    val participant_id: String,
    val participant_name: String,
    val total_badges: Int,
    val total_points: Int,
    val achievements: List<Achievement>,
    val available_badges: List<Badge>
)

data class ParticipantRank(
    val participant_id: String,
    val naam: String,
    val rank: Int,
    val total_score: Int,
    val steps: Int,
    val achievement_points: Int,
    val badge_count: Int,
    val above_me: LeaderboardEntry?,
    val below_me: LeaderboardEntry?
)
```

### 2. API Service

```kotlin
// File: api/GamificationApi.kt

import retrofit2.Response
import retrofit2.http.*

interface GamificationApi {
    
    // Public Endpoints
    @GET("api/leaderboard")
    suspend fun getLeaderboard(
        @Query("page") page: Int = 1,
        @Query("limit") limit: Int = 50,
        @Query("route") route: String? = null,
        @Query("year") year: Int? = null,
        @Query("min_steps") minSteps: Int? = null,
        @Query("top_n") topN: Int? = null
    ): Response<LeaderboardResponse>
    
    @GET("api/badges")
    suspend fun getBadges(
        @Header("Authorization") token: String? = null
    ): Response<List<Badge>>
    
    // Participant Endpoints (Authenticated)
    @GET("api/participant/achievements")
    suspend fun getMyAchievements(
        @Header("Authorization") token: String
    ): Response<ParticipantSummary>
    
    @GET("api/participant/rank")
    suspend fun getMyRank(
        @Header("Authorization") token: String
    ): Response<ParticipantRank>
    
    @GET("api/participant/{id}/achievements")
    suspend fun getParticipantAchievements(
        @Header("Authorization") token: String,
        @Path("id") participantId: String
    ): Response<ParticipantSummary>
}
```

### 3. Repository Pattern

```kotlin
// File: repository/GamificationRepository.kt

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class GamificationRepository(
    private val api: GamificationApi,
    private val authStorage: AuthStorage,
    private val cacheManager: CacheManager
) {
    
    // Leaderboard with caching
    suspend fun getLeaderboard(
        page: Int = 1,
        limit: Int = 50,
        route: String? = null,
        forceRefresh: Boolean = false
    ): Result<LeaderboardResponse> = withContext(Dispatchers.IO) {
        val cacheKey = "leaderboard_${page}_${limit}_${route}"
        
        // Check cache first (5 min TTL)
        if (!forceRefresh) {
            val cached = cacheManager.get<LeaderboardResponse>(cacheKey)
            if (cached != null) {
                return@withContext Result.success(cached)
            }
        }
        
        try {
            val response = api.getLeaderboard(page, limit, route)
            
            if (response.isSuccessful) {
                val leaderboard = response.body()!!
                cacheManager.put(cacheKey, leaderboard, ttlMinutes = 5)
                Result.success(leaderboard)
            } else {
                Result.failure(Exception("Failed to fetch leaderboard: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    // Badges with long-term caching
    suspend fun getBadges(forceRefresh: Boolean = false): Result<List<Badge>> = 
        withContext(Dispatchers.IO) {
            val cacheKey = "badges_all"
            
            if (!forceRefresh) {
                val cached = cacheManager.get<List<Badge>>(cacheKey)
                if (cached != null) {
                    return@withContext Result.success(cached)
                }
            }
            
            try {
                val token = authStorage.getToken()
                val authHeader = token?.let { "Bearer $it" }
                
                val response = api.getBadges(authHeader)
                
                if (response.isSuccessful) {
                    val badges = response.body()!!
                    cacheManager.put(cacheKey, badges, ttlMinutes = 60)
                    Result.success(badges)
                } else {
                    Result.failure(Exception("Failed to fetch badges"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
    
    // My Achievements (no caching - always fresh)
    suspend fun getMyAchievements(): Result<ParticipantSummary> = 
        withContext(Dispatchers.IO) {
            try {
                val token = authStorage.getToken() 
                    ?: return@withContext Result.failure(Exception("Not authenticated"))
                
                val response = api.getMyAchievements("Bearer $token")
                
                if (response.isSuccessful) {
                    Result.success(response.body()!!)
                } else {
                    Result.failure(Exception("Failed to fetch achievements"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
    
    // My Rank with short cache
    suspend fun getMyRank(forceRefresh: Boolean = false): Result<ParticipantRank> = 
        withContext(Dispatchers.IO) {
            val cacheKey = "my_rank"
            
            if (!forceRefresh) {
                val cached = cacheManager.get<ParticipantRank>(cacheKey)
                if (cached != null) {
                    return@withContext Result.success(cached)
                }
            }
            
            try {
                val token = authStorage.getToken()
                    ?: return@withContext Result.failure(Exception("Not authenticated"))
                
                val response = api.getMyRank("Bearer $token")
                
                if (response.isSuccessful) {
                    val rank = response.body()!!
                    cacheManager.put(cacheKey, rank, ttlMinutes = 2)
                    Result.success(rank)
                } else {
                    Result.failure(Exception("Failed to fetch rank"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
```

### 4. Cache Manager

```kotlin
// File: cache/CacheManager.kt

import com.google.gson.Gson
import java.util.concurrent.ConcurrentHashMap

data class CacheEntry<T>(
    val data: T,
    val expiresAt: Long
)

class CacheManager(private val gson: Gson = Gson()) {
    private val cache = ConcurrentHashMap<String, CacheEntry<*>>()
    
    fun <T> put(key: String, data: T, ttlMinutes: Int) {
        val expiresAt = System.currentTimeMillis() + (ttlMinutes * 60 * 1000)
        cache[key] = CacheEntry(data, expiresAt)
    }
    
    @Suppress("UNCHECKED_CAST")
    fun <T> get(key: String): T? {
        val entry = cache[key] as? CacheEntry<T> ?: return null
        
        if (System.currentTimeMillis() > entry.expiresAt) {
            cache.remove(key)
            return null
        }
        
        return entry.data
    }
    
    fun clear() {
        cache.clear()
    }
    
    fun clearKey(key: String) {
        cache.remove(key)
    }
}
```

### 5. UI Components (Jetpack Compose)

```kotlin
// File: ui/LeaderboardScreen.kt

@Composable
fun LeaderboardScreen(
    viewModel: GamificationViewModel = hiltViewModel()
) {
    val leaderboard by viewModel.leaderboard.collectAsState()
    val loading by viewModel.loading.collectAsState()
    
    Column(modifier = Modifier.fillMaxSize()) {
        // Header
        TopAppBar(
            title = { Text("Leaderboard") },
            actions = {
                IconButton(onClick = { viewModel.refresh() }) {
                    Icon(Icons.Default.Refresh, "Refresh")
                }
            }
        )
        
        // Loading State
        if (loading) {
            Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                CircularProgressIndicator()
            }
            return@Column
        }
        
        // Leaderboard List
        LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp)
        ) {
            leaderboard?.entries?.let { entries ->
                itemsIndexed(entries) { index, entry ->
                    LeaderboardEntryCard(entry = entry)
                    if (index < entries.lastIndex) {
                        Spacer(modifier = Modifier.height(8.dp))
                    }
                }
            }
        }
    }
}

@Composable
fun LeaderboardEntryCard(entry: LeaderboardEntry) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Rank Badge
            Box(
                modifier = Modifier
                    .size(40.dp)
                    .background(
                        color = when (entry.rank) {
                            1 -> Color(0xFFFFD700) // Gold
                            2 -> Color(0xFFC0C0C0) // Silver
                            3 -> Color(0xFFCD7F32) // Bronze
                            else -> Color.Gray
                        },
                        shape = CircleShape
                    ),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "#${entry.rank}",
                    color = Color.White,
                    fontWeight = FontWeight.Bold
                )
            }
            
            Spacer(modifier = Modifier.width(12.dp))
            
            // Info
            Column(modifier = Modifier.weight(1f)) {
                Text(entry.naam, fontWeight = FontWeight.Bold)
                Text("${entry.route} • ${entry.steps.toLocaleString()} steps")
            }
            
            // Score
            Column(horizontalAlignment = Alignment.End) {
                Text(
                    text = entry.total_score.toLocaleString(),
                    fontWeight = FontWeight.Bold,
                    fontSize = 18.sp
                )
                Row {
                    Icon(Icons.Default.Star, "", tint = Color(0xFFFFD700), modifier = Modifier.size(16.dp))
                    Text(" ${entry.badge_count}", fontSize = 12.sp)
                }
            }
        }
    }
}

// Extension for number formatting
fun Int.toLocaleString(): String {
    return String.format("%,d", this)
}
```

---

## 🍎 iOS Implementation (Swift)

### 1. Data Models

```swift
// File: Models/Gamification.swift

import Foundation

struct LeaderboardEntry: Codable {
    let id: String
    let naam: String
    let route: String
    let steps: Int
    let achievement_points: Int
    let total_score: Int
    let rank: Int
    let badge_count: Int
    let joined_at: String
}

struct LeaderboardResponse: Codable {
    let entries: [LeaderboardEntry]
    let total_entries: Int
    let current_page: Int
    let total_pages: Int
    let limit: Int
}

struct Badge: Codable {
    let id: String
    let name: String
    let description: String
    let icon_url: String?
    let criteria: BadgeCriteria
    let points: Int
    let is_active: Bool
    let display_order: Int
    let created_at: String
    let updated_at: String
    let earned_count: Int?
    let last_earned_at: String?
    let earned_by_current_user: Bool?
}

struct BadgeCriteria: Codable {
    let min_steps: Int?
    let min_days: Int?
    let consecutive_days: Int?
    let route: String?
    let early_participant: Bool?
    let has_team: Bool?
}

struct Achievement: Codable {
    let id: String
    let participant_id: String
    let earned_at: String
    let badge: Badge
}

struct ParticipantSummary: Codable {
    let participant_id: String
    let participant_name: String
    let total_badges: Int
    let total_points: Int
    let achievements: [Achievement]
    let available_badges: [Badge]
}

struct ParticipantRank: Codable {
    let participant_id: String
    let naam: String
    let rank: Int
    let total_score: Int
    let steps: Int
    let achievement_points: Int
    let badge_count: Int
    let above_me: LeaderboardEntry?
    let below_me: LeaderboardEntry?
}
```

### 2. API Service

```swift
// File: Services/GamificationService.swift

import Foundation

class GamificationService {
    private let baseURL = "https://api.dekoninklijkeloop.nl"
    private let authStorage = AuthStorage()
    private let cacheManager = CacheManager()
    
    // MARK: - Leaderboard
    
    func getLeaderboard(
        page: Int = 1,
        limit: Int = 50,
        route: String? = nil,
        forceRefresh: Bool = false
    ) async throws -> LeaderboardResponse {
        let cacheKey = "leaderboard_\(page)_\(limit)_\(route ?? "")"
        
        // Check cache (5 min TTL)
        if !forceRefresh, let cached: LeaderboardResponse = cacheManager.get(key: cacheKey) {
            return cached
        }
        
        var components = URLComponents(string: "\(baseURL)/api/leaderboard")!
        components.queryItems = [
            URLQueryItem(name: "page", value: String(page)),
            URLQueryItem(name: "limit", value: String(limit))
        ]
        
        if let route = route {
            components.queryItems?.append(URLQueryItem(name: "route", value: route))
        }
        
        let request = URLRequest(url: components.url!)
        let (data, _) = try await URLSession.shared.data(for: request)
        
        let leaderboard = try JSONDecoder().decode(LeaderboardResponse.self, from: data)
        cacheManager.put(key: cacheKey, value: leaderboard, ttlMinutes: 5)
        
        return leaderboard
    }
    
    // MARK: - Badges
    
    func getBadges(forceRefresh: Bool = false) async throws -> [Badge] {
        let cacheKey = "badges_all"
        
        if !forceRefresh, let cached: [Badge] = cache Manager.get(key: cacheKey) {
            return cached
        }
        
        let url = URL(string: "\(baseURL)/api/badges")!
        var request = URLRequest(url: url)
        
        if let token = authStorage.getToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        let (data, _) = try await URLSession.shared.data(for: request)
        let badges = try JSONDecoder().decode([Badge].self, from: data)
        
        cacheManager.put(key: cacheKey, value: badges, ttlMinutes: 60)
        return badges
    }
    
    // MARK: - Achievements
    
    func getMyAchievements() async throws -> ParticipantSummary {
        guard let token = authStorage.getToken() else {
            throw GamificationError.notAuthenticated
        }
        
        let url = URL(string: "\(baseURL)/api/participant/achievements")!
        var request = URLRequest(url: url)
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, _) = try await URLSession.shared.data(for: request)
        return try JSONDecoder().decode(ParticipantSummary.self, from: data)
    }
    
    // MARK: - Rank
    
    func getMyRank(forceRefresh: Bool = false) async throws -> ParticipantRank {
        let cacheKey = "my_rank"
        
        if !forceRefresh, let cached: ParticipantRank = cacheManager.get(key: cacheKey) {
            return cached
        }
        
        guard let token = authStorage.getToken() else {
            throw GamificationError.notAuthenticated
        }
        
        let url = URL(string: "\(baseURL)/api/participant/rank")!
        var request = URLRequest(url: url)
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, _) = try await URLSession.shared.data(for: request)
        let rank = try JSONDecoder().decode(ParticipantRank.self, from: data)
        
        cacheManager.put(key: cacheKey, value: rank, ttlMinutes: 2)
        return rank
    }
}

enum GamificationError: Error {
    case notAuthenticated
    case networkError
}
```

### 3. Cache Manager

```swift
// File: Services/CacheManager.swift

import Foundation

struct CacheEntry<T: Codable>: Codable {
    let data: T
    let expiresAt: Date
}

class CacheManager {
    private let userDefaults = UserDefaults.standard
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()
    
    func put<T: Codable>(key: String, value: T, ttlMinutes: Int) {
        let expiresAt = Date().addingTimeInterval(TimeInterval(ttlMinutes * 60))
        let entry = CacheEntry(data: value, expiresAt: expiresAt)
        
        if let encoded = try? encoder.encode(entry) {
            userDefaults.set(encoded, forKey: "cache_\(key)")
        }
    }
    
    func get<T: Codable>(key: String) -> T? {
        guard let data = userDefaults.data(forKey: "cache_\(key)") else {
            return nil
        }
        
        guard let entry = try? decoder.decode(CacheEntry<T>.self, from: data) else {
            return nil
        }
        
        if Date() > entry.expiresAt {
            clear(key: key)
            return nil
        }
        
        return entry.data
    }
    
    func clear(key: String) {
        userDefaults.removeObject(forKey: "cache_\(key)")
    }
    
    func clearAll() {
        let keys = userDefaults.dictionaryRepresentation().keys
        keys.filter { $0.hasPrefix("cache_") }.forEach { userDefaults.removeObject(forKey: $0) }
    }
}
```

### 4. SwiftUI Views

```swift
// File: Views/LeaderboardView.swift

import SwiftUI

struct LeaderboardView: View {
    @StateObject private var viewModel = LeaderboardViewModel()
    
    var body: some View {
        NavigationView {
            Group {
                if viewModel.loading {
                    ProgressView()
                } else if let error = viewModel.error {
                    ErrorView(error: error) {
                        Task { await viewModel.refresh() }
                    }
                } else {
                    ScrollView {
                        LazyVStack(spacing: 12) {
                            ForEach(viewModel.entries) { entry in
                                LeaderboardEntryCard(entry: entry)
                            }
                        }
                        .padding()
                    }
                }
            }
            .navigationTitle("Leaderboard")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { Task { await viewModel.refresh() } }) {
                        Image(systemName: "arrow.clockwise")
                    }
                }
            }
        }
        .task {
            await viewModel.loadLeaderboard()
        }
    }
}

struct LeaderboardEntryCard: View {
    let entry: LeaderboardEntry
    
    var body: some View {
        HStack(spacing: 12) {
            // Rank Badge
            ZStack {
                Circle()
                    .fill(rankColor)
                    .frame(width: 40, height: 40)
                
                Text("#\(entry.rank)")
                    .font(.caption)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
            }
            
            // Info
            VStack(alignment: .leading, spacing: 4) {
                Text(entry.naam)
                    .font(.headline)
                
                Text("\(entry.route) • \(entry.steps.formatted()) steps")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            
            Spacer()
            
            // Score
            VStack(alignment: .trailing, spacing: 4) {
                Text(entry.total_score.formatted())
                    .font(.title3)
                    .fontWeight(.bold)
                
                HStack(spacing: 2) {
                    Image(systemName: "star.fill")
                        .foregroundColor(.yellow)
                        .font(.caption)
                    Text("\(entry.badge_count)")
                        .font(.caption)
                }
            }
        }
        .padding()
        .background(Color(.systemGray6))
        .cornerRadius(12)
    }
    
    var rankColor: Color {
        switch entry.rank {
        case 1: return Color(hex: "#FFD700") // Gold
        case 2: return Color(hex: "#C0C0C0") // Silver
        case 3: return Color(hex: "#CD7F32") // Bronze
        default: return .gray
        }
    }
}

// Color extension for hex
extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue:  Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
```

---

## ⚛️ React Native Implementation

```typescript
// File: services/gamificationService.ts

import { AuthStorage } from './authStorage';

const BASE_URL = 'https://api.dekoninklijkeloop.nl';

interface LeaderboardEntry {
  id: string;
  naam: string;
  route: string;
  steps: number;
  achievement_points: number;
  total_score: number;
  rank: number;
  badge_count: number;
  joined_at: string;
}

interface LeaderboardResponse {
  entries: LeaderboardEntry[];
  total_entries: number;
  current_page: number;
  total_pages: number;
  limit: number;
}

// Cache implementation
class Cache {
  private cache = new Map<string, { data: any; expiresAt: number }>();

  put(key: string, data: any, ttlMinutes: number) {
    const expiresAt = Date.now() + ttlMinutes * 60 * 1000;
    this.cache.set(key, { data, expiresAt });
  }

  get<T>(key: string): T | null {
    const entry = this.cache.get(key);
    if (!entry) return null;

    if (Date.now() > entry.expiresAt) {
      this.cache.delete(key);
      return null;
    }

    return entry.data as T;
  }

  clear() {
    this.cache.clear();
  }
}

const cache = new Cache();

export class GamificationService {
  // Leaderboard with caching
  static async getLeaderboard(
    page: number = 1,
    limit: number = 50,
    route?: string,
    forceRefresh: boolean = false
  ): Promise<LeaderboardResponse> {
    const cacheKey = `leaderboard_${page}_${limit}_${route || ''}`;

    if (!forceRefresh) {
      const cached = cache.get<LeaderboardResponse>(cacheKey);
      if (cached) return cached;
    }

    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...(route && { route }),
    });

    const response = await fetch(`${BASE_URL}/api/leaderboard?${params}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch leaderboard');
    }

    const data = await response.json();
    cache.put(cacheKey, data, 5); // 5 min cache
    
    return data;
  }

  // Get My Achievements
  static async getMyAchievements(): Promise<any> {
    const token = await AuthStorage.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch(`${BASE_URL}/api/participant/achievements`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch achievements');
    }

    return response.json();
  }
  
  // Get My Rank
  static async getMyRank(forceRefresh: boolean = false): Promise<any> {
    const cacheKey = 'my_rank';

    if (!forceRefresh) {
      const cached = cache.get(cacheKey);
      if (cached) return cached;
    }

    const token = await AuthStorage.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch(`${BASE_URL}/api/participant/rank`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch rank');
    }

    const data = await response.json();
    cache.put(cacheKey, data, 2); // 2 min cache
    
    return data;
  }
}
```

---

## 💙 Flutter Implementation

```dart
// File: lib/services/gamification_service.dart

import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/gamification.dart';
import 'auth_storage.dart';
import 'cache_manager.dart';

class GamificationService {
  static const String baseUrl = 'https://api.dekoninklijkeloop.nl';
  final CacheManager _cache = CacheManager();
  
  // Get Leaderboard
  Future<LeaderboardResponse> getLeaderboard({
    int page = 1,
    int limit = 50,
    String? route,
    bool forceRefresh = false,
  }) async {
    final cacheKey = 'leaderboard_${page}_${limit}_${route ?? ""}';
    
    if (!forceRefresh) {
      final cached = _cache.get<LeaderboardResponse>(cacheKey);
      if (cached != null) return cached;
    }
    
    final queryParams = {
      'page': page.toString(),
      'limit': limit.toString(),
      if (route != null) 'route': route,
    };
    
    final uri = Uri.parse('$baseUrl/api/leaderboard').replace(
      queryParameters: queryParams,
    );
    
    final response = await http.get(uri);
    
    if (response.statusCode != 200) {
      throw Exception('Failed to load leaderboard');
    }
    
    final leaderboard = LeaderboardResponse.fromJson(
      jsonDecode(response.body),
    );
    
    _cache.put(cacheKey, leaderboard, ttlMinutes: 5);
    return leaderboard;
  }
  
  // Get My Achievements
  Future<ParticipantSummary> getMyAchievements() async {
    final token = await AuthStorage.getToken();
    if (token == null) {
      throw Exception('Not authenticated');
    }
    
    final response = await http.get(
      Uri.parse('$baseUrl/api/participant/achievements'),
      headers: {'Authorization': 'Bearer $token'},
    );
    
    if (response.statusCode != 200) {
      throw Exception('Failed to load achievements');
    }
    
    return ParticipantSummary.fromJson(jsonDecode(response.body));
  }
  
  // Get My Rank
  Future<ParticipantRank> getMyRank({bool forceRefresh = false}) async {
    const cacheKey = 'my_rank';
    
    if (!forceRefresh) {
      final cached = _cache.get<ParticipantRank>(cacheKey);
      if (cached != null) return cached;
    }
    
    final token = await AuthStorage.getToken();
    if (token == null) {
      throw Exception('Not authenticated');
    }
    
    final response = await http.get(
      Uri.parse('$baseUrl/api/participant/rank'),
      headers: {'Authorization': 'Bearer $token'},
    );
    
    if (response.statusCode != 200) {
      throw Exception('Failed to load rank');
    }
    
    final rank = ParticipantRank.fromJson(jsonDecode(response.body));
    _cache.put(cacheKey, rank, ttlMinutes: 2);
    
    return rank;
  }
}
```

---

## 💾 Offline Support

### Android (Room Database)

```kotlin
// Cache leaderboard offline
@Entity(tableName = "cached_leaderboard")
data class CachedLeaderboard(
    @PrimaryKey val cacheKey: String,
    val data: String, // JSON
    val cachedAt: Long
)

@Dao
interface LeaderboardDao {
    @Query("SELECT * FROM cached_leaderboard WHERE cacheKey = :key AND cachedAt > :minTime")
    suspend fun getCached(key: String, minTime: Long): CachedLeaderboard?
    
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun cache(leaderboard: CachedLeaderboard)
}
```

### iOS (Core Data)

```swift
// Create CachedData entity in Core Data
class CacheStore {
    func save<T: Codable>(key: String, value: T, ttlMinutes: Int) {
        // Save to Core Data SQLite
    }
    
    func get<T: Codable>(key: String) -> T? {
        // Fetch from Core Data SQLite
    }
}
```

---

## ⚡ Performance Tips

### 1. Pagination Strategy

```typescript
// Load more pattern
const [page, setPage] = useState(1);
const [entries, setEntries] = useState<LeaderboardEntry[]>([]);

const loadMore = async () => {
  const response = await GamificationService.getLeaderboard(page + 1, 20);
  setEntries([...entries, ...response.entries]);
  setPage(page + 1);
};
```

### 2. Pull-to-Refresh

```kotlin
// Android SwipeRefreshLayout
swipeRefresh.setOnRefreshListener {
    viewModel.refresh()
}
```

```swift
// iOS RefreshControl
List {
    ForEach(entries) { entry in
        LeaderboardEntryCard(entry: entry)
    }
}
.refreshable {
    await viewModel.refresh()
}
```

### 3. Image Caching

```kotlin
// Use Coil for badge icons (Android)
AsyncImage(
    model = ImageRequest.Builder(context)
        .data(badge.icon_url)
        .crossfade(true)
        .diskCachePolicy(CachePolicy.ENABLED)
        .memoryCachePolicy(CachePolicy.ENABLED)
        .build(),
    contentDescription = badge.name
)
```

```swift
// Use SDWebImage for badge icons (iOS)
AsyncImage(url: URL(string: badge.icon_url ?? "")) { image in
    image.resizable()
} placeholder: {
    ProgressView()
}
```

---

## 🐛 Troubleshooting

### Common Issues

**1. 401 Unauthorized on participant endpoints**
- Verify token is included in Authorization header
- Check token expiry (refresh if needed)
- Ensure user is logged in

**2. Slow leaderboard loading**
- Reduce limit parameter (default 50 → 20)
- Use pagination instead of loading all at once
- Check cache implementation

**3. Badges not displaying**
- Verify icon_url is valid HTTP/HTTPS URL
- Check CORS if loading from web
- Implement fallback placeholder image

**4. Cache not working**
- Verify TTL is set correctly
- Check cache key consistency
- Clear cache on logout

---

## 📚 Additional Resources

- [Main Gamification API Docs](./GAMIFICATION_API.md)
- [Mobile RBAC Migration Guide](../MOBILE_RBAC_MIGRATION.md)
- [Auth Implementation](../AUTH_AND_RBAC.md)

---

**Version:** 1.0  
**Last Updated:** 2025-11-02  
**API Base:** `https://api.dekoninklijkeloop.nl`  
**Status:** ✅ Mobile Optimized