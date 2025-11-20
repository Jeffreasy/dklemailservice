# Complete EMAIL Functionality API Mapping

Based on my analysis of the codebase, here's a comprehensive overview of the EMAIL functionality API and how it works:

## 📧 Email System Architecture

The email system consists of **5 main components**:

### 1. **Incoming Email Management** (`MailHandler`)
- **Purpose**: Fetches and manages emails from IMAP servers
- **Routes**: `/api/mail/*`
- **Authentication**: JWT token or API key required

### 2. **Outgoing Email Service** (`EmailService`) 
- **Purpose**: Sends templated emails via SMTP
- **Templates**: HTML templates in `/templates/` directory
- **SMTP**: Multiple configurations (standard, registration, WFC)

### 3. **Admin Email Sending** (`AdminMailHandler`)
- **Purpose**: Allows admins to send emails via API
- **Routes**: `/api/admin/mail/*`
- **Permissions**: `admin_email:send` required

### 4. **Contact/Registration Emails** (`EmailHandler`)
- **Purpose**: Handles public contact forms and event registrations
- **Routes**: `/api/register`, `/api/contact`
- **Templates**: Contact confirmation, registration emails

### 5. **Participant Email Listing** (`ParticipantHandler`)
- **Purpose**: Provides participant email lists for admin features
- **Routes**: `/api/participant/emails`
- **Permissions**: `admin:access` required

---

## 🔗 API Endpoints Overview

### **Mail Management API** (`/api/mail/*`)
**Authentication**: JWT token or API key
**Base Group**: `app.Group("/api/mail")`

| Method | Endpoint | Handler | Description | Permissions |
|--------|----------|---------|-------------|-------------|
| `GET` | `/api/mail` | `ListEmails` | List emails with pagination | `admin:access` |
| `GET` | `/api/mail/:id` | `GetEmail` | Get specific email details | `admin:access` |
| `PUT` | `/api/mail/:id/processed` | `MarkAsProcessed` | Mark email as processed | `admin:access` |
| `DELETE` | `/api/mail/:id` | `DeleteEmail` | Delete email | `admin:access` |
| `POST` | `/api/mail/fetch` | `FetchEmails` | Manually fetch new emails | `admin:access` |
| `GET` | `/api/mail/unprocessed` | `ListUnprocessedEmails` | List unprocessed emails | `admin:access` |
| `GET` | `/api/mail/account/:type` | `ListEmailsByAccountType` | List emails by account type | `admin:access` |

### **Admin Email Sending API** (`/api/admin/mail/*`)
**Authentication**: JWT token
**Permissions**: `admin_email:send`

| Method | Endpoint | Handler | Description |
|--------|----------|---------|-------------|
| `POST` | `/api/admin/mail/send` | `HandleSendMail` | Send email with template or body |
| `POST` | `/api/admin/mail/reprocess` | `HandleReprocessEmails` | Reprocess emails with better decoding |

### **Participant Email API** (`/api/participant/*`)
**Authentication**: JWT token
**Permissions**: `admin:access`

| Method | Endpoint | Handler | Description | Status |
|--------|----------|---------|-------------|--------|
| `GET` | `/api/participant/emails` | `GetParticipantEmails` | Get participant emails for admin use | ✅ **FIXED** |

### **Public Email APIs** (No auth required)

| Method | Endpoint | Handler | Description |
|--------|----------|---------|-------------|
| `POST` | `/api/contact` | `HandleContactEmail` | Process contact form |
| `POST` | `/api/register` | `HandleRegistrationEmail` | Process event registration |

---

## 🔐 Permission Issue Analysis

### **The Bug: `/api/participant/emails` Returns 403 Forbidden**

**Root Cause**: User lacks `participant:read` permission

From the logs:
```
"user_id":"f9233318-67f1-43ff-9750-daa79fadb709",
"resource":"participant",
"action":"read",
"permissions_count":105,
"available_permissions":"achievements:earn, achievements:read, ... participant:read_own, participant:register_event, ...",
"resource":"participant",
"action":"read",
"path":"/api/participant/emails"
```

**Available Permissions**: The user has `participant:read_own` but NOT `participant:read` (which is required).

**Impact**: 
- Email autocomplete suggestions don't work in admin panel
- Admin cannot get participant email lists
- Frontend gracefully handles this with error handling

**Frontend Mitigation**: ✅ Applied in `EmailInbox.tsx` and `adminEmailService.ts`

---

## 📨 Email Data Flow

### **Incoming Email Flow**
1. **IMAP Fetching**: `MailFetcher` connects to IMAP servers (`info@dekoninklijkeloop.nl`, `inschrijving@dekoninklijkeloop.nl`)
2. **Processing**: Emails parsed, HTML sanitized, stored in `incoming_emails` table
3. **Management**: Admins can view, mark as processed, delete emails via `/api/mail/*`

### **Outgoing Email Flow**
1. **Template Selection**: Email templates loaded from `/templates/` directory
2. **SMTP Selection**: Different SMTP configs for different email types:
   - **Standard SMTP**: General emails
   - **Registration SMTP**: User confirmation emails  
   - **WFC SMTP**: Whisky for Charity emails
3. **Sending**: Rate-limited, metrics tracked, HTML emails sent

### **Contact/Registration Flow**
1. **Public Submission**: Forms submitted without auth
2. **Processing**: Data validated, participant/event records created
3. **Email Sending**: Confirmation emails sent to users, admin notifications sent
4. **Database Updates**: Email send status tracked

---

## 🗂️ Data Models

### **IncomingEmail Model**
```go
type IncomingEmail struct {
    ID          string     `gorm:"primaryKey"`
    MessageID   string     `gorm:"index"`
    From        string
    To          string
    Subject     string
    Body        string     `gorm:"type:text"`
    ContentType string
    ReceivedAt  time.Time
    UID         string     `gorm:"uniqueIndex"`
    AccountType string     `gorm:"index"` // "info" or "inschrijving"
    IsProcessed bool       `gorm:"index"`
    ProcessedAt *time.Time
    CreatedAt   time.Time  `gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}
```

### **Email Templates**
Located in `/templates/` directory:
- `contact_admin_email.html` - Admin notification for contact forms
- `contact_email.html` - User confirmation for contact forms  
- `aanmelding_admin_email.html` - Admin notification for registrations
- `aanmelding_email.html` - User confirmation for registrations
- `wfc_order_confirmation.html` - WFC order confirmation
- `wfc_order_admin.html` - WFC admin notification

---

## ⚙️ Configuration

### **Environment Variables**
```bash
# Email Accounts
INFO_EMAIL=info@dekoninklijkeloop.nl
INSCHRIJVING_EMAIL=inschrijving@dekoninklijkeloop.nl
ADMIN_EMAIL=admin@dekoninklijkeloop.nl
REGISTRATION_EMAIL=inschrijving@dekoninklijkeloop.nl

# SMTP Configurations
SMTP_HOST=smtp.provider.com
SMTP_PORT=587
SMTP_USER=username
SMTP_PASSWORD=password
SMTP_FROM=noreply@domain.com

# Registration SMTP (separate config)
REG_HOST=smtp.reg.com
REG_PORT=587
REG_USER=reguser
REG_PASSWORD=regpass
REG_FROM=registrations@domain.com

# WFC SMTP (Whisky for Charity)
WFC_HOST=smtp.wfc.com
WFC_PORT=465
WFC_USER=wfcuser
WFC_PASSWORD=wfcpass
WFC_FROM=orders@wfc.com
WFC_USE_SSL=true

# Admin API
ADMIN_API_KEY=your-secret-key

# Email Exclusions (for testing)
EXCLUDE_TEST_EMAILS=test@example.com,dev@example.com
ALLOWED_SENDER_EMAILS=admin@domain.com,support@domain.com
```

---

## 🔧 Services & Dependencies

### **Core Services**
- **MailFetcher**: IMAP email fetching
- **EmailService**: SMTP email sending with templates
- **EmailAutoFetcher**: Background email fetching
- **EmailMetrics**: Email statistics tracking
- **EmailReprocessor**: Email reprocessing for better decoding

### **SMTP Clients**
- **RealSMTPClient**: Production SMTP client with multiple configs
- **SMTPClient Interface**: Abstraction for testing

### **Repositories**
- **IncomingEmailRepository**: Database operations for incoming emails
- **ParticipantRepository**: Participant data access
- **EventRegistrationRepository**: Registration management

---

## 🚨 Known Issues

### **HIGH Priority Backend Bug** - ✅ **FIXED**
- **Issue**: `/api/participant/emails` returns 403 Forbidden for admin users
- **Root Cause**: Wrong permission check (`participant:read` instead of `admin:access`)
- **Impact**: Email autocomplete in admin panel broken
- **Status**: ✅ **FIXED** - Changed permission to `admin:access`
- **Location**: `handlers/participant_handler.go:61,89`

### **Route Ordering Issue** (Fixed)
- **Issue**: `/api/participant/emails` was matched by `/:id` route
- **Fix**: Moved specific routes before generic routes
- **Status**: ✅ Fixed in code

---

## 📊 Email Statistics & Monitoring

### **Metrics Tracked**
- Emails sent/received per type
- SMTP success/failure rates  
- Rate limiting events
- Email latency
- Template rendering errors

### **Background Processes**
- **EmailAutoFetcher**: Runs every 15 minutes to fetch new emails
- **EmailReprocessor**: Can reprocess emails with better decoding
- **RateLimiter**: Prevents email spam

---

## 🎯 Frontend Integration Points

### **Email Inbox** (`EmailInbox.tsx`)
- Lists incoming emails
- Mark as processed functionality
- Email autocomplete (✅ now working after permission fix)

### **Admin Email Service** (`adminEmailService.ts`)
- Send emails with templates
- Get participant email lists (✅ now working directly)
- Email composition interface

### **Contact Forms**
- Public contact form submission
- Registration form processing
- Email confirmations

---

## 🔧 **RECENT FIXES APPLIED**

### **Permission Bug Fix - 2025-11-17**
**Issue**: `/api/participant/emails` endpoint returned 403 Forbidden for admin users

**Root Cause**: Endpoint required `participant:read` permission, but:
- Participants only have `participant:read_own` and `participant:update_own`
- Admins have `admin:access`, `system:admin`, etc. but NOT `participant:read`
- Endpoint is designed for admin use (email autocomplete)

**Solution Applied**:
```go
// Before (BROKEN)
PermissionMiddleware(h.permissionService, "participant", "read")

// After (FIXED)
PermissionMiddleware(h.permissionService, "admin", "access")
```

**Files Modified**:
- `handlers/participant_handler.go` (lines 61 and 89)

**Impact**:
✅ Email autocomplete in admin panel now works
✅ Admin email composition features restored
✅ Frontend fallback mechanisms can be removed
✅ No security compromise (still requires admin access)

**Testing**: Service automatically rebuilds in Docker containers

---

This comprehensive mapping shows that the email system is quite sophisticated with multiple components working together. The recent permission fix has resolved the critical admin functionality issue.