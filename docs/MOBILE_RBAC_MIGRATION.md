# 📱 Mobile App RBAC Migration Guide

> **Version:** 1.0 | **Target:** Android & iOS Apps | **Backend:** V1.48.0+ | **Last Updated:** 2025-11-02

Complete migratie gids voor mobiele applicaties (Android & iOS) naar het nieuwe RBAC systeem van DKL Email Service.

---

## 📋 Table of Contents

1. [Breaking Changes Overview](#-breaking-changes-overview)
2. [API Response Changes](#-api-response-changes)
3. [Android Implementation](#-android-implementation)
4. [iOS Implementation](#-ios-implementation)
5. [React Native Implementation](#-react-native-implementation)
6. [Flutter Implementation](#-flutter-implementation)
7. [Migration Checklist](#-migration-checklist)
8. [Testing Guide](#-testing-guide)
9. [Troubleshooting](#-troubleshooting)

---

## ⚠️ Breaking Changes Overview

### Wat is Er Veranderd?

**VOOR (Legacy System):**
```json
{
  "user": {
    "id": "uuid",
    "naam": "Admin",
    "email": "admin@dekoninklijkeloop.nl",
    "rol": "admin"  // ❌ VERWIJDERD - Single string
  }
}
```

**NA (RBAC System V1.22+):**
```json
{
  "user": {
    "id": "uuid",
    "naam": "Admin",
    "email": "admin@dekoninklijkeloop.nl",
    "roles": [  // ✅ NIEUW - Array van role objecten
      {
        "id": "role-uuid",
        "name": "admin",
        "description": "Administrator"
      }
    ],
    "permissions": [  // ✅ NIEUW - Granulaire permissions
      {"resource": "admin", "action": "access"},
      {"resource": "contact", "action": "read"},
      {"resource": "contact", "action": "write"}
    ]
  }
}
```

### Waarom Deze Change?

- **Enterprise RBAC:** Granulaire permissions in plaats van simpele roles
- **Multi-role Support:** Gebruikers kunnen meerdere rollen hebben
- **Feature Flags:** Fine-grained feature toegang via permissions
- **Security:** Backend-driven authorization (geen hardcoded roles)

---

## 📡 API Response Changes

### Login Endpoint

**Endpoint:** `POST /api/auth/login`

**Request:**
```json
{
  "email": "admin@dekoninklijkeloop.nl",
  "wachtwoord": "password"
}
```

**Response (Nieuw):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "base64-encoded-refresh-token",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@dekoninklijkeloop.nl",
    "naam": "Admin Gebruiker",
    "is_actief": true,
    "roles": [
      {
        "id": "role-uuid",
        "name": "admin",
        "description": "Administrator met volledige toegang"
      }
    ],
    "permissions": [
      {"resource": "admin", "action": "access"},
      {"resource": "contact", "action": "read"},
      {"resource": "contact", "action": "write"},
      {"resource": "contact", "action": "delete"},
      {"resource": "user", "action": "read"},
      {"resource": "user", "action": "manage_roles"}
    ]
  }
}
```

### Profile Endpoint

**Endpoint:** `GET /api/auth/profile`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "naam": "Admin Gebruiker",
  "email": "admin@dekoninklijkeloop.nl",
  "is_actief": true,
  "laatste_login": "2025-11-02T20:00:00Z",
  "created_at": "2025-01-01T00:00:00Z",
  "roles": [
    {
      "id": "role-uuid",
      "name": "admin",
      "description": "Administrator",
      "assigned_at": "2025-01-01T00:00:00Z",
      "is_active": true
    }
  ],
  "permissions": [
    {"resource": "admin", "action": "access"},
    {"resource": "contact", "action": "read"}
  ]
}
```

### Token Refresh

**Endpoint:** `POST /api/auth/refresh`

**Request:**
```json
{
  "refresh_token": "base64-encoded-refresh-token"
}
```

**Response:**
```json
{
  "success": true,
  "token": "new-jwt-token",
  "refresh_token": "new-refresh-token"
}
```

---

## 📱 Android Implementation

### 1. Data Models (Kotlin)

```kotlin
// File: models/Auth.kt

data class LoginRequest(
    val email: String,
    val wachtwoord: String
)

data class LoginResponse(
    val success: Boolean,
    val token: String,
    val refresh_token: String,
    val user: User
)

data class User(
    val id: String,
    val naam: String,
    val email: String,
    val is_actief: Boolean,
    val roles: List<Role>,
    val permissions: List<Permission>,
    val laatste_login: String? = null,
    val created_at: String? = null
)

data class Role(
    val id: String,
    val name: String,
    val description: String,
    val assigned_at: String? = null,
    val is_active: Boolean? = null
)

data class Permission(
    val resource: String,
    val action: String
)

data class RefreshRequest(
    val refresh_token: String
)

data class RefreshResponse(
    val success: Boolean,
    val token: String,
    val refresh_token: String
)
```

### 2. Storage Manager (Kotlin)

```kotlin
// File: storage/AuthStorage.kt

import android.content.Context
import android.content.SharedPreferences
import com.google.gson.Gson

class AuthStorage(context: Context) {
    private val prefs: SharedPreferences = 
        context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
    private val gson = Gson()

    // Token Management
    fun saveToken(token: String) {
        prefs.edit().putString(KEY_TOKEN, token).apply()
    }

    fun getToken(): String? {
        return prefs.getString(KEY_TOKEN, null)
    }

    fun saveRefreshToken(refreshToken: String) {
        prefs.edit().putString(KEY_REFRESH_TOKEN, refreshToken).apply()
    }

    fun getRefreshToken(): String? {
        return prefs.getString(KEY_REFRESH_TOKEN, null)
    }

    // User Data Management
    fun saveUser(user: User) {
        val userJson = gson.toJson(user)
        prefs.edit().putString(KEY_USER, userJson).apply()
    }

    fun getUser(): User? {
        val userJson = prefs.getString(KEY_USER, null) ?: return null
        return try {
            gson.fromJson(userJson, User::class.java)
        } catch (e: Exception) {
            null
        }
    }

    // Role & Permission Helpers
    fun getUserRoles(): List<Role> {
        return getUser()?.roles ?: emptyList()
    }

    fun getUserPermissions(): List<Permission> {
        return getUser()?.permissions ?: emptyList()
    }

    fun getPrimaryRole(): String {
        return getUserRoles().firstOrNull()?.name ?: "user"
    }

    // Permission Checking
    fun hasPermission(resource: String, action: String): Boolean {
        return getUserPermissions().any { 
            it.resource == resource && it.action == action 
        }
    }

    fun hasAnyPermission(vararg checks: Pair<String, String>): Boolean {
        val permissions = getUserPermissions()
        return checks.any { (resource, action) ->
            permissions.any { it.resource == resource && it.action == action }
        }
    }

    fun hasRole(roleName: String): Boolean {
        return getUserRoles().any { it.name == roleName }
    }

    // Clear All
    fun clear() {
        prefs.edit().clear().apply()
    }

    companion object {
        private const val KEY_TOKEN = "auth_token"
        private const val KEY_REFRESH_TOKEN = "refresh_token"
        private const val KEY_USER = "user_data"
    }
}
```

### 3. API Service (Kotlin + Retrofit)

```kotlin
// File: api/AuthApi.kt

import retrofit2.Response
import retrofit2.http.*

interface AuthApi {
    @POST("api/auth/login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>

    @GET("api/auth/profile")
    suspend fun getProfile(
        @Header("Authorization") token: String
    ): Response<User>

    @POST("api/auth/refresh")
    suspend fun refreshToken(@Body request: RefreshRequest): Response<RefreshResponse>

    @POST("api/auth/logout")
    suspend fun logout(
        @Header("Authorization") token: String
    ): Response<Unit>
}
```

### 4. Auth Repository (Kotlin)

```kotlin
// File: repository/AuthRepository.kt

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class AuthRepository(
    private val authApi: AuthApi,
    private val authStorage: AuthStorage
) {
    suspend fun login(email: String, password: String): Result<User> = withContext(Dispatchers.IO) {
        try {
            val response = authApi.login(LoginRequest(email, password))
            
            if (response.isSuccessful) {
                val loginResponse = response.body()!!
                
                // Save tokens
                authStorage.saveToken(loginResponse.token)
                authStorage.saveRefreshToken(loginResponse.refresh_token)
                
                // Save user data (includes roles & permissions)
                authStorage.saveUser(loginResponse.user)
                
                Result.success(loginResponse.user)
            } else {
                Result.failure(Exception("Login failed: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun refreshToken(): Result<String> = withContext(Dispatchers.IO) {
        try {
            val refreshToken = authStorage.getRefreshToken()
                ?: return@withContext Result.failure(Exception("No refresh token"))

            val response = authApi.refreshToken(RefreshRequest(refreshToken))
            
            if (response.isSuccessful) {
                val refreshResponse = response.body()!!
                
                // Save new tokens
                authStorage.saveToken(refreshResponse.token)
                authStorage.saveRefreshToken(refreshResponse.refresh_token)
                
                Result.success(refreshResponse.token)
            } else {
                Result.failure(Exception("Token refresh failed"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProfile(): Result<User> = withContext(Dispatchers.IO) {
        try {
            val token = authStorage.getToken()
                ?: return@withContext Result.failure(Exception("Not authenticated"))

            val response = authApi.getProfile("Bearer $token")
            
            if (response.isSuccessful) {
                val user = response.body()!!
                authStorage.saveUser(user)
                Result.success(user)
            } else {
                Result.failure(Exception("Profile fetch failed"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    fun isAuthenticated(): Boolean {
        return authStorage.getToken() != null
    }

    suspend fun logout() {
        try {
            val token = authStorage.getToken()
            if (token != null) {
                authApi.logout("Bearer $token")
            }
        } finally {
            authStorage.clear()
        }
    }
}
```

### 5. Permission Helper Extensions (Kotlin)

```kotlin
// File: extensions/PermissionExtensions.kt

// Check single permission
fun User.hasPermission(resource: String, action: String): Boolean {
    return permissions.any { it.resource == resource && it.action == action }
}

// Check any permission
fun User.hasAnyPermission(vararg checks: Pair<String, String>): Boolean {
    return checks.any { (resource, action) -> hasPermission(resource, action) }
}

// Check all permissions
fun User.hasAllPermissions(vararg checks: Pair<String, String>): Boolean {
    return checks.all { (resource, action) -> hasPermission(resource, action) }
}

// Check role
fun User.hasRole(roleName: String): Boolean {
    return roles.any { it.name == roleName }
}

// Quick admin check
fun User.isAdmin(): Boolean {
    return hasPermission("admin", "access")
}

// Quick staff check
fun User.isStaff(): Boolean {
    return hasPermission("staff", "access")
}

// Usage examples:
// if (user.hasPermission("contact", "write")) { ... }
// if (user.hasAnyPermission("contact" to "read", "contact" to "write")) { ... }
// if (user.isAdmin()) { ... }
```

---

## 🍎 iOS Implementation

### 1. Data Models (Swift)

```swift
// File: Models/Auth.swift

import Foundation

struct LoginRequest: Codable {
    let email: String
    let wachtwoord: String
}

struct LoginResponse: Codable {
    let success: Bool
    let token: String
    let refresh_token: String
    let user: User
}

struct User: Codable {
    let id: String
    let naam: String
    let email: String
    let is_actief: Bool
    let roles: [Role]
    let permissions: [Permission]
    let laatste_login: String?
    let created_at: String?
}

struct Role: Codable {
    let id: String
    let name: String
    let description: String
    let assigned_at: String?
    let is_active: Bool?
}

struct Permission: Codable {
    let resource: String
    let action: String
}

struct RefreshRequest: Codable {
    let refresh_token: String
}

struct RefreshResponse: Codable {
    let success: Bool
    let token: String
    let refresh_token: String
}
```

### 2. Storage Manager (Swift)

```swift
// File: Storage/AuthStorage.swift

import Foundation

class AuthStorage {
    private let defaults = UserDefaults.standard
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()
    
    private enum Keys {
        static let token = "auth_token"
        static let refreshToken = "refresh_token"
        static let user = "user_data"
    }
    
    // MARK: - Token Management
    
    func saveToken(_ token: String) {
        defaults.set(token, forKey: Keys.token)
    }
    
    func getToken() -> String? {
        return defaults.string(forKey: Keys.token)
    }
    
    func saveRefreshToken(_ refreshToken: String) {
        defaults.set(refreshToken, forKey: Keys.refreshToken)
    }
    
    func getRefreshToken() -> String? {
        return defaults.string(forKey: Keys.refreshToken)
    }
    
    // MARK: - User Data Management
    
    func saveUser(_ user: User) {
        if let encoded = try? encoder.encode(user) {
            defaults.set(encoded, forKey: Keys.user)
        }
    }
    
    func getUser() -> User? {
        guard let data = defaults.data(forKey: Keys.user) else { return nil }
        return try? decoder.decode(User.self, from: data)
    }
    
    // MARK: - Role & Permission Helpers
    
    func getUserRoles() -> [Role] {
        return getUser()?.roles ?? []
    }
    
    func getUserPermissions() -> [Permission] {
        return getUser()?.permissions ?? []
    }
    
    func getPrimaryRole() -> String {
        return getUserRoles().first?.name ?? "user"
    }
    
    // MARK: - Permission Checking
    
    func hasPermission(resource: String, action: String) -> Bool {
        return getUserPermissions().contains { 
            $0.resource == resource && $0.action == action 
        }
    }
    
    func hasAnyPermission(_ checks: [(resource: String, action: String)]) -> Bool {
        let permissions = getUserPermissions()
        return checks.contains { resource, action in
            permissions.contains { $0.resource == resource && $0.action == action }
        }
    }
    
    func hasRole(_ roleName: String) -> Bool {
        return getUserRoles().contains { $0.name == roleName }
    }
    
    // MARK: - Clear
    
    func clear() {
        defaults.removeObject(forKey: Keys.token)
        defaults.removeObject(forKey: Keys.refreshToken)
        defaults.removeObject(forKey: Keys.user)
    }
}
```

### 3. API Service (Swift)

```swift
// File: Services/AuthService.swift

import Foundation

class AuthService {
    private let baseURL = "https://api.dekoninklijkeloop.nl"
    private let storage = AuthStorage()
    
    // MARK: - Login
    
    func login(email: String, password: String) async throws -> User {
        let url = URL(string: "\(baseURL)/api/auth/login")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let loginRequest = LoginRequest(email: email, wachtwoord: password)
        request.httpBody = try JSONEncoder().encode(loginRequest)
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw AuthError.loginFailed
        }
        
        let loginResponse = try JSONDecoder().decode(LoginResponse.self, from: data)
        
        // Save tokens and user
        storage.saveToken(loginResponse.token)
        storage.saveRefreshToken(loginResponse.refresh_token)
        storage.saveUser(loginResponse.user)
        
        return loginResponse.user
    }
    
    // MARK: - Refresh Token
    
    func refreshToken() async throws -> String {
        guard let refreshToken = storage.getRefreshToken() else {
            throw AuthError.noRefreshToken
        }
        
        let url = URL(string: "\(baseURL)/api/auth/refresh")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let refreshRequest = RefreshRequest(refresh_token: refreshToken)
        request.httpBody = try JSONEncoder().encode(refreshRequest)
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw AuthError.tokenRefreshFailed
        }
        
        let refreshResponse = try JSONDecoder().decode(RefreshResponse.self, from: data)
        
        // Save new tokens
        storage.saveToken(refreshResponse.token)
        storage.saveRefreshToken(refreshResponse.refresh_token)
        
        return refreshResponse.token
    }
    
    // MARK: - Get Profile
    
    func getProfile() async throws -> User {
        guard let token = storage.getToken() else {
            throw AuthError.notAuthenticated
        }
        
        let url = URL(string: "\(baseURL)/api/auth/profile")!
        var request = URLRequest(url: url)
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw AuthError.profileFetchFailed
        }
        
        let user = try JSONDecoder().decode(User.self, from: data)
        storage.saveUser(user)
        
        return user
    }
    
    // MARK: - Logout
    
    func logout() async {
        if let token = storage.getToken() {
            let url = URL(string: "\(baseURL)/api/auth/logout")!
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
            
            try? await URLSession.shared.data(for: request)
        }
        
        storage.clear()
    }
    
    // MARK: - Status
    
    func isAuthenticated() -> Bool {
        return storage.getToken() != nil
    }
}

enum AuthError: Error {
    case loginFailed
    case noRefreshToken
    case tokenRefreshFailed
    case notAuthenticated
    case profileFetchFailed
}
```

### 4. Permission Extensions (Swift)

```swift
// File: Extensions/User+Permissions.swift

extension User {
    // Check single permission
    func hasPermission(resource: String, action: String) -> Bool {
        return permissions.contains { $0.resource == resource && $0.action == action }
    }
    
    // Check any permission
    func hasAnyPermission(_ checks: [(resource: String, action: String)]) -> Bool {
        return checks.contains { resource, action in
            hasPermission(resource: resource, action: action)
        }
    }
    
    // Check all permissions
    func hasAllPermissions(_ checks: [(resource: String, action: String)]) -> Bool {
        return checks.allSatisfy { resource, action in
            hasPermission(resource: resource, action: action)
        }
    }
    
    // Check role
    func hasRole(_ roleName: String) -> Bool {
        return roles.contains { $0.name == roleName }
    }
    
    // Quick checks
    var isAdmin: Bool {
        return hasPermission(resource: "admin", action: "access")
    }
    
    var isStaff: Bool {
        return hasPermission(resource: "staff", action: "access")
    }
}

// Usage examples:
// if user.hasPermission(resource: "contact", action: "write") { ... }
// if user.hasAnyPermission([("contact", "read"), ("contact", "write")]) { ... }
// if user.isAdmin { ... }
```

---

## ⚛️ React Native Implementation

### 1. Data Types (TypeScript)

```typescript
// File: types/auth.ts

export interface LoginRequest {
  email: string;
  wachtwoord: string;
}

export interface LoginResponse {
  success: boolean;
  token: string;
  refresh_token: string;
  user: User;
}

export interface User {
  id: string;
  naam: string;
  email: string;
  is_actief: boolean;
  roles: Role[];
  permissions: Permission[];
  laatste_login?: string;
  created_at?: string;
}

export interface Role {
  id: string;
  name: string;
  description: string;
  assigned_at?: string;
  is_active?: boolean;
}

export interface Permission {
  resource: string;
  action: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface RefreshResponse {
  success: boolean;
  token: string;
  refresh_token: string;
}
```

### 2. Storage Service (TypeScript)

```typescript
// File: services/authStorage.ts

import AsyncStorage from '@react-native-async-storage/async-storage';
import { User, Role, Permission } from '../types/auth';

const KEYS = {
  TOKEN: 'auth_token',
  REFRESH_TOKEN: 'refresh_token',
  USER: 'user_data',
};

export class AuthStorage {
  // Token Management
  static async saveToken(token: string): Promise<void> {
    await AsyncStorage.setItem(KEYS.TOKEN, token);
  }

  static async getToken(): Promise<string | null> {
    return await AsyncStorage.getItem(KEYS.TOKEN);
  }

  static async saveRefreshToken(refreshToken: string): Promise<void> {
    await AsyncStorage.setItem(KEYS.REFRESH_TOKEN, refreshToken);
  }

  static async getRefreshToken(): Promise<string | null> {
    return await AsyncStorage.getItem(KEYS.REFRESH_TOKEN);
  }

  // User Data Management
  static async saveUser(user: User): Promise<void> {
    await AsyncStorage.setItem(KEYS.USER, JSON.stringify(user));
  }

  static async getUser(): Promise<User | null> {
    const userJson = await AsyncStorage.getItem(KEYS.USER);
    if (!userJson) return null;
    
    try {
      return JSON.parse(userJson) as User;
    } catch (error) {
      console.error('Error parsing user data:', error);
      return null;
    }
  }

  // Role & Permission Helpers
  static async getUserRoles(): Promise<Role[]> {
    const user = await this.getUser();
    return user?.roles ?? [];
  }

  static async getUserPermissions(): Promise<Permission[]> {
    const user = await this.getUser();
    return user?.permissions ?? [];
  }

  static async getPrimaryRole(): Promise<string> {
    const roles = await this.getUserRoles();
    return roles[0]?.name ?? 'user';
  }

  // Permission Checking
  static async hasPermission(resource: string, action: string): Promise<boolean> {
    const permissions = await this.getUserPermissions();
    return permissions.some(p => p.resource === resource && p.action === action);
  }

  static async hasAnyPermission(...checks: [string, string][]): Promise<boolean> {
    const permissions = await this.getUserPermissions();
    return checks.some(([resource, action]) =>
      permissions.some(p => p.resource === resource && p.action === action)
    );
  }

  static async hasRole(roleName: string): Promise<boolean> {
    const roles = await this.getUserRoles();
    return roles.some(r => r.name === roleName);
  }

  // Clear All
  static async clear(): Promise<void> {
    await AsyncStorage.multiRemove([KEYS.TOKEN, KEYS.REFRESH_TOKEN, KEYS.USER]);
  }
}
```

### 3. API Service (TypeScript)

```typescript
// File: services/authService.ts

import { AuthStorage } from './authStorage';
import {
  LoginRequest,
  LoginResponse,
  RefreshRequest,
  RefreshResponse,
  User,
} from '../types/auth';

const BASE_URL = 'https://api.dekoninklijkeloop.nl';

export class AuthService {
  // Login
  static async login(email: string, password: string): Promise<User> {
    const response = await fetch(`${BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, wachtwoord: password } as LoginRequest),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data: LoginResponse = await response.json();

    // Save tokens and user
    await AuthStorage.saveToken(data.token);
    await AuthStorage.saveRefreshToken(data.refresh_token);
    await AuthStorage.saveUser(data.user);

    return data.user;
  }

  // Refresh Token
  static async refreshToken(): Promise<string> {
    const refreshToken = await AuthStorage.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await fetch(`${BASE_URL}/api/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken } as RefreshRequest),
    });

    if (!response.ok) {
      throw new Error('Token refresh failed');
    }

    const data: RefreshResponse = await response.json();

    // Save new tokens
    await AuthStorage.saveToken(data.token);
    await AuthStorage.saveRefreshToken(data.refresh_token);

    return data.token;
  }

  // Get Profile
  static async getProfile(): Promise<User> {
    const token = await AuthStorage.getToken();
    if (!token) {
      throw new Error('Not authenticated');
    }

    const response = await fetch(`${BASE_URL}/api/auth/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      throw new Error('Profile fetch failed');
    }

    const user: User = await response.json();
    await AuthStorage.saveUser(user);

    return user;
  }

  // Logout
  static async logout(): Promise<void> {
    const token = await AuthStorage.getToken();
    
    if (token) {
      try {
        await fetch(`${BASE_URL}/api/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });
      } catch (error) {
        console.error('Logout API call failed:', error);
      }
    }

    await AuthStorage.clear();
  }

  // Is Authenticated
  static async isAuthenticated(): Promise<boolean> {
    const token = await AuthStorage.getToken();
    return token !== null;
  }
}
```

### 4. Permission Hooks (TypeScript)

```typescript
// File: hooks/usePermissions.ts

import { useState, useEffect } from 'react';
import { AuthStorage } from '../services/authStorage';
import { Permission } from '../types/auth';

export const usePermissions = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPermissions();
  }, []);

  const loadPermissions = async () => {
    try {
      const perms = await AuthStorage.getUserPermissions();
      setPermissions(perms);
    } finally {
      setLoading(false);
    }
  };

  const hasPermission = (resource: string, action: string): boolean => {
    return permissions.some(p => p.resource === resource && p.action === action);
  };

  const hasAnyPermission = (...checks: [string, string][]): boolean => {
    return checks.some(([resource, action]) => hasPermission(resource, action));
  };

  const hasAllPermissions = (...checks: [string, string][]): boolean => {
    return checks.every(([resource, action]) => hasPermission(resource, action));
  };

  return {
    permissions,
    loading,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    refresh: loadPermissions,
  };
};

// Usage:
// const { hasPermission } = usePermissions();
// if (hasPermission('contact', 'write')) { ... }
```

---

## 💙 Flutter Implementation

### 1. Data Models (Dart)

```dart
// File: lib/models/auth.dart

import 'dart:convert';

class LoginRequest {
  final String email;
  final String wachtwoord;

  LoginRequest({required this.email, required this.wachtwoord});

  Map<String, dynamic> toJson() => {
    'email': email,
    'wachtwoord': wachtwoord,
  };
}

class LoginResponse {
  final bool success;
  final String token;
  final String refreshToken;
  final User user;

  LoginResponse({
    required this.success,
    required this.token,
    required this.refreshToken,
    required this.user,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) => LoginResponse(
    success: json['success'],
    token: json['token'],
    refreshToken: json['refresh_token'],
    user: User.fromJson(json['user']),
  );
}

class User {
  final String id;
  final String naam;
  final String email;
  final bool isActief;
  final List<Role> roles;
  final List<Permission> permissions;
  final String? laatsteLogin;
  final String? createdAt;

  User({
    required this.id,
    required this.naam,
    required this.email,
    required this.isActief,
    required this.roles,
    required this.permissions,
    this.laatsteLogin,
    this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) => User(
    id: json['id'],
    naam: json['naam'],
    email: json['email'],
    isActief: json['is_actief'],
    roles: (json['roles'] as List).map((r) => Role.fromJson(r)).toList(),
    permissions: (json['permissions'] as List).map((p) => Permission.fromJson(p)).toList(),
    laatsteLogin: json['laatste_login'],
    createdAt: json['created_at'],
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'naam': naam,
    'email': email,
    'is_actief': isActief,
    'roles': roles.map((r) => r.toJson()).toList(),
    'permissions': permissions.map((p) => p.toJson()).toList(),
    'laatste_login': laatsteLogin,
    'created_at': createdAt,
  };
}

class Role {
  final String id;
  final String name;
  final String description;
  final String? assignedAt;
  final bool? isActive;

  Role({
    required this.id,
    required this.name,
    required this.description,
    this.assignedAt,
    this.isActive,
  });

  factory Role.fromJson(Map<String, dynamic> json) => Role(
    id: json['id'],
    name: json['name'],
    description: json['description'],
    assignedAt: json['assigned_at'],
    isActive: json['is_active'],
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'assigned_at': assignedAt,
    'is_active': isActive,
  };
}

class Permission {
  final String resource;
  final String action;

  Permission({required this.resource, required this.action});

  factory Permission.fromJson(Map<String, dynamic> json) => Permission(
    resource: json['resource'],
    action: json['action'],
  );

  Map<String, dynamic> toJson() => {
    'resource': resource,
    'action': action,
  };
}
```

### 2. Storage Service (Dart)

```dart
// File: lib/services/auth_storage.dart

import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/auth.dart';

class AuthStorage {
  static const String _keyToken = 'auth_token';
  static const String _keyRefreshToken = 'refresh_token';
  static const String _keyUser = 'user_data';

  // Token Management
  static Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_keyToken, token);
  }

  static Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_keyToken);
  }

  static Future<void> saveRefreshToken(String refreshToken) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_keyRefreshToken, refreshToken);
  }

  static Future<String?> getRefreshToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_keyRefreshToken);
  }

  // User Data Management
  static Future<void> saveUser(User user) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_keyUser, jsonEncode(user.toJson()));
  }

  static Future<User?> getUser() async {
    final prefs = await SharedPreferences.getInstance();
    final userJson = prefs.getString(_keyUser);
    
    if (userJson == null) return null;
    
    try {
      return User.fromJson(jsonDecode(userJson));
    } catch (e) {
      print('Error parsing user data: $e');
      return null;
    }
  }

  // Role & Permission Helpers
  static Future<List<Role>> getUserRoles() async {
    final user = await getUser();
    return user?.roles ?? [];
  }

  static Future<List<Permission>> getUserPermissions() async {
    final user = await getUser();
    return user?.permissions ?? [];
  }

  static Future<String> getPrimaryRole() async {
    final roles = await getUserRoles();
    return roles.isNotEmpty ? roles.first.name : 'user';
  }

  // Permission Checking
  static Future<bool> hasPermission(String resource, String action) async {
    final permissions = await getUserPermissions();
    return permissions.any((p) => p.resource == resource && p.action == action);
  }

  static Future<bool> hasAnyPermission(List<(String, String)> checks) async {
    final permissions = await getUserPermissions();
    return checks.any((check) =>
      permissions.any((p) => p.resource == check.$1 && p.action == check.$2)
    );
  }

  static Future<bool> hasRole(String roleName) async {
    final roles = await getUserRoles();
    return roles.any((r) => r.name == roleName);
  }

  // Clear All
  static Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_keyToken);
    await prefs.remove(_keyRefreshToken);
    await prefs.remove(_keyUser);
  }
}
```

### 3. API Service (Dart)

```dart
// File: lib/services/auth_service.dart

import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/auth.dart';
import 'auth_storage.dart';

class AuthService {
  static const String baseUrl = 'https://api.dekoninklijkeloop.nl';

  // Login
  static Future<User> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/auth/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(LoginRequest(email: email, wachtwoord: password).toJson()),
    );

    if (response.statusCode != 200) {
      throw Exception('Login failed: ${response.statusCode}');
    }

    final loginResponse = LoginResponse.fromJson(jsonDecode(response.body));

    // Save tokens and user
    await AuthStorage.saveToken(loginResponse.token);
    await AuthStorage.saveRefreshToken(loginResponse.refreshToken);
    await AuthStorage.saveUser(loginResponse.user);

    return loginResponse.user;
  }

  // Refresh Token
  static Future<String> refreshToken() async {
    final refreshToken = await AuthStorage.getRefreshToken();
    if (refreshToken == null) {
      throw Exception('No refresh token available');
    }

    final response = await http.post(
      Uri.parse('$baseUrl/api/auth/refresh'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'refresh_token': refreshToken}),
    );

    if (response.statusCode != 200) {
      throw Exception('Token refresh failed');
    }

    final data = jsonDecode(response.body);

    // Save new tokens
    await AuthStorage.saveToken(data['token']);
    await AuthStorage.saveRefreshToken(data['refresh_token']);

    return data['token'];
  }

  // Get Profile
  static Future<User> getProfile() async {
    final token = await AuthStorage.getToken();
    if (token == null) {
      throw Exception('Not authenticated');
    }

    final response = await http.get(
      Uri.parse('$baseUrl/api/auth/profile'),
      headers: {'Authorization': 'Bearer $token'},
    );

    if (response.statusCode != 200) {
      throw Exception('Profile fetch failed');
    }

    final user = User.fromJson(jsonDecode(response.body));
    await AuthStorage.saveUser(user);

    return user;
  }

  // Logout
  static Future<void> logout() async {
    final token = await AuthStorage.getToken();
    
    if (token != null) {
      try {
        await http.post(
          Uri.parse('$baseUrl/api/auth/logout'),
          headers: {'Authorization': 'Bearer $token'},
        );
      } catch (e) {
        print('Logout API call failed: $e');
      }
    }

    await AuthStorage.clear();
  }

  // Is Authenticated
  static Future<bool> isAuthenticated() async {
    final token = await AuthStorage.getToken();
    return token != null;
  }
}
```

---

## ✅ Migration Checklist

### Phase 1: Pre-Migration
- [ ] Backup huidige app versie
- [ ] Test backend API endpoints met Postman/cURL
- [ ] Verify backend versie is V1.48.0+
- [ ] Document bestaande functionaliteit

### Phase 2: Implementation
- [ ] Update data models voor `roles` en `permissions` arrays
- [ ] Verwijder oude `rol` (single string) referenties
- [ ] Implementeer nieuwe storage layer met RBAC support
- [ ] Update API service calls
- [ ] Implementeer permission checking helpers
- [ ] Update UI permission checks

### Phase 3: Testing
- [ ] Test login flow met verschillende roles
- [ ] Test token refresh mechanisme
- [ ] Test permission checks in UI
- [ ] Test offline storage persistence
- [ ] Test logout en data clearing
- [ ] Test edge cases (no permissions, multiple roles)

### Phase 4: Deployment
- [ ] Deploy naar test environment
- [ ] Run smoke tests
- [ ] Deploy naar production
- [ ] Monitor error logs
- [ ] User feedback verzamelen

---

##  Testing Guide

### Test Scenarios

#### 1. Login Test
```
Given: Valid credentials
When: User logs in
Then: 
  - Tokens are saved
  - User object contains roles array
  - User object contains permissions array
  - No 'rol' field present
```

#### 2. Permission Check Test
```
Given: User with 'contact:read' permission
When: hasPermission('contact', 'read') called
Then: Returns true

Given: User without permission
When: hasPermission('contact', 'delete') called
Then: Returns false
```

#### 3. Multi-Role Test
```
Given: User with ['admin', 'staff'] roles
When: Checking roles
Then: 
  - getUserRoles() returns both roles
  - hasRole('admin') returns true
  - hasRole('staff') returns true
```

#### 4. Token Refresh Test
```
Given: Valid refresh token
When: Token expires after 20 minutes
Then: 
  - New access token received
  - New refresh token received
  - Old tokens invalidated
```

### Test Users

Test met deze users in je database:

| Email | Role | Permissions | Use Case |
|-------|------|-------------|----------|
| admin@dekoninklijkeloop.nl | admin | ALL (58 permissions) | Full access test |
| staff@dekoninklijkeloop.nl | staff | Read-only | Limited access test |
| user@dekoninklijkeloop.nl | user | Chat only | Basic user test |

---

## 🐛 Troubleshooting

### AsyncStorage Error: "passing null/undefined"

**Problem:**
```
passed key: userRole
passed value: undefined
```

**Cause:** Old code trying to save non-existent `rol` field

**Fix:**
```typescript
// OLD (❌ WRONG)
await AsyncStorage.setItem('userRole', user.rol);

// NEW (✅ CORRECT)
const primaryRole = user.roles?.[0]?.name ?? 'user';
await AsyncStorage.setItem('userRole', primaryRole);
```

### Empty Permissions Array

**Problem:** User has roles but `permissions` array is empty

**Diagnosis:**
```sql
-- Check database
SELECT u.email, r.name as role, p.resource, p.action
FROM gebruikers u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.email = 'your@email.com';
```

**Fix:** Assign permissions to roles in database or admin panel

### Token Not Refreshing

**Problem:** App logs out after 20 minutes

**Check:**
1. Refresh token saved correctly?
2. Backend refresh endpoint working?
3. Auto-refresh timer implemented?

**Fix:** Implement auto-refresh 5 minutes before token expiry

---

## 📚 Additional Resources

- [AUTH_AND_RBAC.md](./AUTH_AND_RBAC.md) - Complete backend RBAC documentation
- [DATABASE_REFERENCE.md](./DATABASE_REFERENCE.md) - Database schema
- [FRONTEND_INTEGRATION.md](./FRONTEND_INTEGRATION.md) - Web frontend guide

---

## 🆘 Support

**Problemen tijdens migratie?**

1. Check backend logs: `docker logs dkl-backend`
2. Verify API response format met Postman
3. Check AsyncStorage/SharedPreferences data
4. Test met verschillende user roles

**Common Issues:**
- ❌ `undefined rol` → Use `roles` array instead
- ❌ Empty permissions → Check database role_permissions
- ❌ 401 errors → Check token format and expiry
- ❌ 403 errors → User has no required permission

---

**Version:** 1.0  
**Last Updated:** 2025-11-02  
**Backend Compatibility:** V1.48.0+  
**Status:** ✅ Production Ready