# V30 Frontend Implementation Guide - Duaal Registratiesysteem

**Versie:** 1.0  
**Datum:** 2025-11-10  
**Backend Versie:** V30  
**Frontend Locatie:** `C:\Users\jeffrey\Desktop\Githubmains\DKL25\src\pages\Aanmelden\`

---

## 🎯 Doel

Implementeer een duaal registratiesysteem waarbij gebruikers tijdens registratie kunnen kiezen tussen:
1. **Volledig Account** - Met DKL Step App toegang
2. **Tijdelijke Registratie** - Alleen voor het event van 2026

## 📦 Benodigde Bestanden

### Te Wijzigen
- ✏️ `src/pages/Aanmelden/components/FormContainer.tsx`
- ✏️ `src/pages/Aanmelden/types/schema.ts`
- ✏️ `src/pages/Aanmelden/components/SuccessMessage.tsx`

### Te Maken (Nieuw)
- ✨ `src/pages/Upgrade/UpgradeAccount.tsx`
- ✨ `src/pages/Aanmelden/components/AccountTypeSelector.tsx`
- ✨ `src/pages/Aanmelden/components/PasswordField.tsx`

---

## 🔧 Stap-voor-Stap Implementatie

### STAP 1: Update Type Definities

**Bestand:** `src/pages/Aanmelden/types/schema.ts`

```typescript
import { z } from 'zod';

export const RegistrationSchema = z.object({
  // Bestaande velden
  naam: z.string().min(1, 'Naam is verplicht'),
  email: z.string().email('Ongeldig e-mailadres'),
  telefoon: z.string().optional(),
  rol: z.enum(['Deelnemer', 'Begeleider', 'Vrijwilliger'], {
    errorMap: () => ({ message: 'Selecteer een rol' })
  }),
  afstand: z.enum(['2.5 KM', '6 KM', '10 KM', '15 KM'], {
    errorMap: () => ({ message: 'Selecteer een afstand' })
  }),
  ondersteuning: z.enum(['Ja', 'Nee', 'Anders'], {
    errorMap: () => ({ message: 'Geef aan of je ondersteuning nodig hebt' })
  }),
  bijzonderheden: z.string().optional(),
  terms: z.boolean().refine(val => val === true, {
    message: 'Je moet akkoord gaan met de voorwaarden'
  }),
  
  // 🆕 NIEUW: Account type velden
  want_account: z.boolean({
    errorMap: () => ({ message: 'Maak een keuze voor account type' })
  }),
  wachtwoord: z.string().min(8, 'Wachtwoord moet minimaal 8 karakters zijn').optional(),
  
  test_mode: z.boolean().optional()
})
// 🆕 NIEUW: Custom validaties
.refine(
  (data) => !data.want_account || (data.wachtwoord && data.wachtwoord.length >= 8),
  {
    message: "Wachtwoord is verplicht voor een volledig account (min. 8 karakters)",
    path: ["wachtwoord"]
  }
)
.refine(
  (data) => {
    if (data.rol === 'Begeleider' || data.rol === 'Vrijwilliger') {
      return data.telefoon && data.telefoon.length > 0;
    }
    return true;
  },
  {
    message: "Telefoonnummer is verplicht voor begeleiders en vrijwilligers",
    path: ["telefoon"]
  }
)
.refine(
  (data) => {
    if (data.ondersteuning === 'Ja' || data.ondersteuning === 'Anders') {
      return data.bijzonderheden && data.bijzonderheden.length > 0;
    }
    return true;
  },
  {
    message: "Bijzonderheden zijn verplicht als je ondersteuning nodig hebt",
    path: ["bijzonderheden"]
  }
);

export type RegistrationFormData = z.infer<typeof RegistrationSchema>;

// 🆕 NIEUW: Response types
export interface PublicRegistrationResponse {
  success: boolean;
  message: string;
  participant_id: string;
  registration_id?: string;
  account_type: 'full' | 'temporary';
  has_app_access: boolean;
  gebruiker_id?: string;
  event_name?: string;
  event_date?: string;
}

export interface UpgradeResponse {
  success: boolean;
  message: string;
  gebruiker_id: string;
  participant_id: string;
  has_app_access: boolean;
}
```

---

### STAP 2: Maak AccountTypeSelector Component

**Nieuw Bestand:** `src/pages/Aanmelden/components/AccountTypeSelector.tsx`

```typescript
import { memo } from 'react';
import { cc, cn, colors } from '@/styles/shared';

interface AccountTypeSelectorProps {
  value: boolean | null;
  onChange: (wantAccount: boolean) => void;
  error?: string;
}

export const AccountTypeSelector = memo<AccountTypeSelectorProps>(({ 
  value, 
  onChange, 
  error 
}) => {
  return (
    <section className="space-y-6" aria-labelledby="account-type-heading">
      <h2 id="account-type-heading" className={cn(cc.text.h2, 'text-gray-900 pb-4 relative', cc.typography.heading)}>
        Wil je een account aanmaken?
      </h2>
      
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
        <p className="text-sm text-blue-900">
          <strong>💡 Wat is het verschil?</strong><br/>
          Met een <strong>volledig account</strong> krijg je toegang tot de DKL Step App en kun je je voortgang bijhouden.
          Met <strong>alleen registreren</strong> schrijf je je in voor het evenement zonder app toegang (je kunt later altijd upgraden).
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4" role="radiogroup" aria-label="Kies account type">
        {/* Full Account Optie */}
        <label className="relative cursor-pointer group">
          <input
            type="radio"
            name="account_choice"
            checked={value === true}
            onChange={() => onChange(true)}
            className="peer sr-only"
            aria-label="Ja, maak een volledig account aan"
          />
          <div className={cn(
            'p-6 rounded-xl border-2 transition-all min-h-[240px]',
            'flex flex-col items-start',
            'border-gray-200 bg-white hover:shadow-md',
            'peer-checked:border-primary peer-checked:bg-primary/5 peer-checked:shadow-lg'
          )}>
            <span className="text-5xl mb-3 group-hover:scale-110 transition-transform">📱</span>
            <h3 className="font-bold text-lg mb-2 text-gray-900">
              Ja, maak een account aan
            </h3>
            <p className="text-sm text-gray-600 mb-4">
              Krijg volledige toegang tot alle DKL features
            </p>
            <ul className="text-sm text-gray-700 space-y-2 flex-grow">
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Toegang tot DKL Step App</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Stappen tracking & voortgang</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Badges & achievements</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Community features</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Permanent account (ook voor volgende jaren)</span>
              </li>
            </ul>
            <div className="mt-4 w-full">
              <span className={cn(
                'inline-block px-3 py-1 rounded-full text-xs font-semibold',
                value === true ? 'bg-primary text-white' : 'bg-gray-100 text-gray-600'
              )}>
                Wachtwoord vereist
              </span>
            </div>
          </div>
        </label>

        {/* Temporary Registration Optie */}
        <label className="relative cursor-pointer group">
          <input
            type="radio"
            name="account_choice"
            checked={value === false}
            onChange={() => onChange(false)}
            className="peer sr-only"
            aria-label="Nee, alleen registreren voor dit evenement"
          />
          <div className={cn(
            'p-6 rounded-xl border-2 transition-all min-h-[240px]',
            'flex flex-col items-start',
            'border-gray-200 bg-white hover:shadow-md',
            'peer-checked:border-primary peer-checked:bg-primary/5 peer-checked:shadow-lg'
          )}>
            <span className="text-5xl mb-3 group-hover:scale-110 transition-transform">📝</span>
            <h3 className="font-bold text-lg mb-2 text-gray-900">
              Nee, alleen registreren
            </h3>
            <p className="text-sm text-gray-600 mb-4">
              Snelle registratie voor het evenement van 2026
            </p>
            <ul className="text-sm text-gray-700 space-y-2 flex-grow">
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Registratie voor evenement 2026</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Geen account nodig</span>
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2 text-lg">✓</span>
                <span>Sneller klaar (geen wachtwoord)</span>
              </li>
              <li className="flex items-start">
                <span className="text-orange-500 mr-2 text-lg">⚠</span>
                <span>Geen app toegang</span>
              </li>
              <li className="flex items-start">
                <span className="text-blue-500 mr-2 text-lg">💡</span>
                <span>Later upgraden mogelijk!</span>
              </li>
            </ul>
            <div className="mt-4 w-full">
              <span className={cn(
                'inline-block px-3 py-1 rounded-full text-xs font-semibold',
                value === false ? 'bg-primary text-white' : 'bg-gray-100 text-gray-600'
              )}>
                Geen wachtwoord nodig
              </span>
            </div>
          </div>
        </label>
      </div>

      {error && (
        <p className={cn(cc.form.errorMessage, 'text-center')} role="alert">
          {error}
        </p>
      )}
    </section>
  );
});

AccountTypeSelector.displayName = 'AccountTypeSelector';
```

---

### STAP 3: Maak PasswordField Component

**Nieuw Bestand:** `src/pages/Aanmelden/components/PasswordField.tsx`

```typescript
import { memo, useState } from 'react';
import { UseFormRegister, FieldErrors } from 'react-hook-form';
import { RegistrationFormData } from '../types/schema';
import { cc, cn } from '@/styles/shared';

interface PasswordFieldProps {
  register: UseFormRegister<RegistrationFormData>;
  errors: FieldErrors<RegistrationFormData>;
}

export const PasswordField = memo<PasswordFieldProps>(({ register, errors }) => {
  const [showPassword, setShowPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState<'weak' | 'medium' | 'strong'>('weak');

  const checkPasswordStrength = (password: string) => {
    if (password.length < 8) return 'weak';
    if (password.length < 12) return 'medium';
    return 'strong';
  };

  return (
    <section 
      className="space-y-6 transition-all duration-300 ease-in-out" 
      aria-labelledby="password-heading"
    >
      <h2 id="password-heading" className={cn(cc.text.h2, 'text-gray-900 pb-4 relative', cc.typography.heading)}>
        Kies je wachtwoord
      </h2>

      <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
        <p className="text-sm text-amber-900">
          <strong>🔒 Wachtwoord voor je account</strong><br/>
          Je gebruikt dit wachtwoord om in te loggen in de DKL Step App en op de website.
        </p>
      </div>

      <div className="space-y-2">
        <label htmlFor="wachtwoord" className={cc.form.label}>
          Wachtwoord <span className="text-red-500">*</span>
        </label>
        
        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            id="wachtwoord"
            className={cn(
              'w-full px-4 py-3 pr-12 rounded-xl border-2',
              'focus:outline-none focus:ring-2 focus:ring-primary/20',
              'text-gray-900 placeholder-gray-400 bg-white',
              errors.wachtwoord ? 'border-red-500' : 'border-gray-200 focus:border-primary'
            )}
            placeholder="Minimaal 8 karakters"
            {...register('wachtwoord', {
              onChange: (e) => {
                setPasswordStrength(checkPasswordStrength(e.target.value));
              }
            })}
          />
          
          {/* Toggle visibility button */}
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
            aria-label={showPassword ? "Verberg wachtwoord" : "Toon wachtwoord"}
          >
            {showPassword ? '👁️' : '👁️‍🗨️'}
          </button>
        </div>

        {/* Password strength indicator */}
        <div className="flex gap-1 h-1">
          <div className={cn(
            'flex-1 rounded transition-colors',
            passwordStrength === 'weak' ? 'bg-red-500' :
            passwordStrength === 'medium' ? 'bg-amber-500' :
            passwordStrength === 'strong' ? 'bg-green-500' : 'bg-gray-200'
          )}/>
          <div className={cn(
            'flex-1 rounded transition-colors',
            passwordStrength === 'medium' || passwordStrength === 'strong' ? 'bg-amber-500' :
            passwordStrength === 'strong' ? 'bg-green-500' : 'bg-gray-200'
          )}/>
          <div className={cn(
            'flex-1 rounded transition-colors',
            passwordStrength === 'strong' ? 'bg-green-500' : 'bg-gray-200'
          )}/>
        </div>
        
        <p className="text-xs text-gray-500">
          {passwordStrength === 'weak' && '🔴 Zwak - Gebruik meer karakters'}
          {passwordStrength === 'medium' && '🟡 Gemiddeld - Voeg speciale tekens toe voor extra veiligheid'}
          {passwordStrength === 'strong' && '🟢 Sterk wachtwoord!'}
        </p>

        {errors.wachtwoord && (
          <p className={cn(cc.form.errorMessage)}>{errors.wachtwoord.message}</p>
        )}

        <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 mt-2">
          <p className="text-xs text-gray-600">
            <strong>💡 Tips voor een sterk wachtwoord:</strong><br/>
            • Gebruik minimaal 8 karakters<br/>
            • Combineer letters, cijfers en speciale tekens<br/>
            • Gebruik geen persoonlijke informatie
          </p>
        </div>
      </div>
    </section>
  );
});

PasswordField.displayName = 'PasswordField';
```

---

### STAP 4: Update FormContainer.tsx

**Bestand:** `src/pages/Aanmelden/components/FormContainer.tsx`

#### 4A. Imports Toevoegen
```typescript
import { AccountTypeSelector } from './AccountTypeSelector';
import { PasswordField } from './PasswordField';
import type { PublicRegistrationResponse } from '../types/schema';
```

#### 4B. Nieuwe State Toevoegen
```typescript
const FormContainer: React.FC<{ onSuccess: (data: RegistrationFormData, response: PublicRegistrationResponse) => void }> = memo(({
  onSuccess
}) => {
  // Bestaande state...
  
  // 🆕 NIEUW: Account type state
  const [accountChoice, setAccountChoice] = useState<boolean | null>(null);
  const [showPasswordField, setShowPasswordField] = useState(false);
  
  // Bestaande form setup...
  const { register, handleSubmit, setValue, control, formState: { errors } } = useForm<RegistrationFormData>({
    resolver: zodResolver(RegistrationSchema),
    defaultValues: {
      ondersteuning: 'Nee' as const,
      bijzonderheden: '',
      want_account: false, // 🆕 NIEUW default
    }
  });
  
  // Watch account choice
  const wantAccount = useWatch({
    control,
    name: 'want_account'
  });
  
  // 🆕 NIEUW: Sync password field visibility
  useEffect(() => {
    setShowPasswordField(wantAccount === true);
    if (wantAccount === false) {
      setValue('wachtwoord', undefined);
    }
  }, [wantAccount, setValue]);
```

#### 4C. Account Type Selector Toevoegen (na Contactgegevens sectie)
```typescript
{/* Account Type Keuze - NIEUW: Moet NA contactgegevens, VOOR rol */}
<AccountTypeSelector
  value={accountChoice}
  onChange={(choice) => {
    setAccountChoice(choice);
    setValue('want_account', choice);
    logEvent('registration', 'select_account_type', choice ? 'full' : 'temporary');
  }}
  error={errors.want_account?.message}
/>
```

#### 4D. Password Field Toevoegen (na Account Type, VOOR Afstand)
```typescript
{/* Wachtwoord sectie - Conditioneel, alleen voor full accounts */}
{showPasswordField && (
  <div className="animate-fadeIn">
    <PasswordField 
      register={register} 
      errors={errors}
    />
  </div>
)}
```

#### 4E. Update onSubmit Functie
```typescript
const onSubmit = async (data: RegistrationFormData) => {
  const startTime = performance.now();

  try {
    logEvent('registration', 'form_submit_attempt', 
      `${data.rol}_${data.afstand}_${data.want_account ? 'full' : 'temp'}`);

    setIsSubmitting(true);
    setSubmitError(null);
    
    const validatedData = validateForm(data);

    // 🆕 NIEUW: Gebruik nieuwe publieke endpoint
    const response = await apiClient.post<PublicRegistrationResponse>(
      '/api/public/aanmelden',
      {
        naam: validatedData.naam,
        email: validatedData.email,
        telefoon: validatedData.telefoon || undefined,
        rol: validatedData.rol,
        afstand: validatedData.afstand,
        ondersteuning: validatedData.ondersteuning,
        bijzonderheden: validatedData.bijzonderheden || '',
        want_account: validatedData.want_account,
        wachtwoord: validatedData.wachtwoord || undefined,
        terms: validatedData.terms,
        test_mode: false
      }
    );

    // Track success met account type
    const endTime = performance.now();
    const duration = Math.round(endTime - startTime);
    logEvent('registration', 'registration_complete', 
      `${validatedData.rol}_${response.account_type}_duration:${duration}ms`);

    // Ga door naar success pagina met response data
    onSuccess(validatedData, response);
    
  } catch (error) {
    const axiosError = error as { response?: { status: number; data?: { error?: string, code?: string } } };
    
    // Handle specific error codes
    if (axiosError?.response?.data?.code === 'ALREADY_REGISTERED') {
      setSubmitError('Je bent al ingeschreven voor dit jaar met dit e-mailadres.');
    } else if (axiosError?.response?.data?.code === 'EMAIL_EXISTS') {
      setSubmitError('Er bestaat al een account met dit e-mailadres. Probeer in te loggen of gebruik een ander e-mailadres.');
    } else if (axiosError?.response?.data?.code === 'VALIDATION_ERROR') {
      setSubmitError(axiosError.response.data.error || 'De ingevoerde gegevens zijn ongeldig.');
    } else {
      setSubmitError('Er ging iets mis bij je aanmelding. Probeer het later opnieuw.');
    }
    
    logEvent('registration', 'form_submit_failure', 
      axiosError?.response?.data?.code || 'unknown_error');
  } finally {
    setIsSubmitting(false);
  }
};
```

---

### STAP 5: Update SuccessMessage.tsx

**Bestand:** `src/pages/Aanmelden/components/SuccessMessage.tsx`

#### 5A. Props Interface Uitbreiden
```typescript
interface SuccessMessageProps {
  data: RegistrationFormData;
  response?: PublicRegistrationResponse; // 🆕 NIEUW: Response data van backend
}

export const SuccessMessage: React.FC<SuccessMessageProps> = memo(({ 
  data, 
  response 
}) => {
  // Bestaande state...
  
  // 🆕 NIEUW: Bepaal of dit een full account is
  const isFullAccount = response?.account_type === 'full';
  const hasAppAccess = response?.has_app_access ?? false;
```

#### 5B. Conditionele Content
```typescript
  return (
    <>
      {showConfetti && <CSSConfetti />}
      
      <article className={cn(cc.container.base, 'py-12 sm:py-16')}>
        <div className={cn('max-w-2xl mx-auto bg-white rounded-xl overflow-hidden', cc.shadow.lg)}>
          
          {/* Header - Aangepast per account type */}
          <header className={cn(colors.primary.bg, 'p-8 text-center')}>
            <div className="mb-4">
              <img src={logoUrl} alt="DKL Logo" className="h-24 mx-auto"/>
            </div>
            <h1 className={cn(cc.text.h2, 'text-white mb-2')}>
              {isFullAccount 
                ? '🎉 Welkom bij De Koninklijke Loop!' 
                : '✅ Bedankt voor je aanmelding!'}
            </h1>
            <p className={cn(cc.text.h5, 'text-white/90')}>
              {isFullAccount
                ? `Je account is aangemaakt en je hebt nu volledige toegang tot de DKL Step App!`
                : `Je bent ingeschreven voor het evenement van 2026.`}
            </p>
          </header>

          {/* App Download Sectie - Alleen voor full accounts */}
          {isFullAccount && hasAppAccess && (
            <section className="p-8 bg-gradient-to-r from-orange-50 to-amber-50 border-y border-orange-100">
              <div className="text-center">
                <h3 className={cn(cc.text.h4, 'text-gray-900 mb-4')}>
                  📱 Download de DKL Step App
                </h3>
                <p className="text-gray-600 mb-6">
                  Begin vandaag nog met stappen verzamelen en verdien rewards!
                </p>
                <div className="flex justify-center gap-4">
                  <a
                    href="https://apps.apple.com/app/dkl-step"
                    target="_blank"
                    rel="noopener noreferrer"
                    className={cn(
                      'inline-flex items-center gap-2 px-6 py-3 bg-black text-white rounded-lg',
                      'hover:bg-gray-800 transition-colors'
                    )}
                  >
                    🍎 App Store
                  </a>
                  <a
                    href="https://play.google.com/store/apps/details?id=nl.dekoninklijkeloop.step"
                    target="_blank"
                    rel="noopener noreferrer"
                    className={cn(
                      'inline-flex items-center gap-2 px-6 py-3 bg-black text-white rounded-lg',
                      'hover:bg-gray-800 transition-colors'
                    )}
                  >
                    📱 Google Play
                  </a>
                </div>
                
                <div className="mt-6 bg-white border border-gray-200 rounded-lg p-4">
                  <p className="text-sm font-semibold text-gray-900 mb-2">
                    🔐 Je inloggegevens:
                  </p>
                  <p className="text-sm text-gray-600">
                    Email: <span className="font-mono">{data.email}</span><br/>
                    Wachtwoord: Het wachtwoord dat je zojuist hebt ingesteld
                  </p>
                </div>
              </div>
            </section>
          )}

          {/* Upgrade CTA - Alleen voor temporary accounts */}
          {!isFullAccount && (
            <section className="p-8 bg-gradient-to-r from-yellow-50 to-amber-50 border-y border-yellow-100">
              <div className="text-center">
                <h3 className={cn(cc.text.h4, 'text-gray-900 mb-4')}>
                  🚀 Wil je meer uit je deelname halen?
                </h3>
                <p className="text-gray-600 mb-4">
                  Upgrade naar een volledig account en krijg toegang tot de DKL Step App!
                </p>
                <ul className="text-left inline-block mb-6 space-y-2">
                  <li className="flex items-center text-sm text-gray-700">
                    <span className="text-green-500 mr-2">✓</span>
                    Track je dagelijkse stappen
                  </li>
                  <li className="flex items-center text-sm text-gray-700">
                    <span className="text-green-500 mr-2">✓</span>
                    Verdien badges en achievements
                  </li>
                  <li className="flex items-center text-sm text-gray-700">
                    <span className="text-green-500 mr-2">✓</span>
                    Join de DKL community
                  </li>
                </ul>
                <br/>
                <a
                  href="/upgrade"
                  className={cn(
                    'inline-flex items-center px-6 py-3 text-white',
                    colors.primary.bg,
                    'rounded-lg hover:bg-primary-dark transition-colors'
                  )}
                >
                  Upgrade nu naar volledig account →
                </a>
              </div>
            </section>
          )}

          {/* Rest van bestaande content... */}
        </div>
      </article>
    </>
  );
});
```

---

### STAP 6: Maak Upgrade Pagina

**Nieuw Bestand:** `src/pages/Upgrade/UpgradeAccount.tsx`

```typescript
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { apiClient } from '@/services/api/apiClient';
import { cc, cn, colors } from '@/styles/shared';
import type { UpgradeResponse } from '../Aanmelden/types/schema';

export const UpgradeAccount: React.FC = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleUpgrade = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Valideer wachtwoorden matchen
    if (password !== confirmPassword) {
      setError('Wachtwoorden komen niet overeen');
      return;
    }

    // Valideer wachtwoord lengte
    if (password.length < 8) {
      setError('Wachtwoord moet minimaal 8 karakters bevatten');
      return;
    }

    try {
      setIsSubmitting(true);

      const response = await apiClient.post<UpgradeResponse>(
        '/api/public/upgrade-to-full-account',
        {
          email,
          wachtwoord: password
        }
      );

      if (response.success) {
        toast.success('🎉 Je account is geüpgraded! Je kunt nu inloggen in de app.');
        
        // Redirect naar app download of login pagina
        setTimeout(() => {
          navigate('/download-app');
        }, 2000);
      }

    } catch (error: any) {
      const errorCode = error?.response?.data?.code;
      const errorMessage = error?.response?.data?.error;

      if (errorCode === 'NO_TEMPORARY_ACCOUNT') {
        setError('Geen tijdelijke registratie gevonden met dit e-mailadres. Heb je je al ingeschreven?');
      } else if (errorCode === 'ALREADY_FULL_ACCOUNT') {
        setError('Dit account is al een volledig account. Je kunt inloggen in de app!');
      } else if (errorCode === 'EMAIL_EXISTS') {
        setError('Er bestaat al een volledig account met dit e-mailadres.');
      } else if (errorCode === 'WEAK_PASSWORD') {
        setError('Wachtwoord moet minimaal 8 karakters bevatten.');
      } else {
        setError(errorMessage || 'Upgrade mislukt. Probeer het later opnieuw.');
      }

      toast.error('Upgrade mislukt');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className={cn(cc.container.base, 'py-12')}>
      <div className="max-w-md mx-auto">
        <div className={cn('bg-white rounded-xl p-8', cc.shadow.lg)}>
          <div className="text-center mb-8">
            <h1 className={cn(cc.text.h2, 'text-gray-900 mb-2')}>
              🚀 Upgrade je Account
            </h1>
            <p className="text-gray-600">
              Krijg toegang tot de DKL Step App en alle features
            </p>
          </div>

          {/* Benefits */}
          <div className="bg-orange-50 border border-orange-200 rounded-lg p-4 mb-6">
            <h3 className="font-semibold text-gray-900 mb-3">
              ✨ Na upgrade krijg je:
            </h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li className="flex items-start">
                <span className="text-green-500 mr-2">✓</span>
                Volledige toegang tot de DKL Step App
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2">✓</span>
                Stappen tracking en voortgang
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2">✓</span>
                Badges en achievements
              </li>
              <li className="flex items-start">
                <span className="text-green-500 mr-2">✓</span>
                Community features
              </li>
            </ul>
          </div>

          <form onSubmit={handleUpgrade} className="space-y-6">
            <div className="space-y-2">
              <label htmlFor="email" className={cc.form.label}>
                E-mailadres <span className="text-red-500">*</span>
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className={cn(
                  'w-full px-4 py-3 rounded-xl border-2',
                  'focus:outline-none focus:ring-2 focus:ring-primary/20',
                  'border-gray-200 focus:border-primary'
                )}
                placeholder="Het e-mailadres waarmee je je hebt geregistreerd"
                required
              />
              <p className="text-xs text-gray-500">
                Gebruik het e-mailadres waarmee je je eerder hebt geregistreerd
              </p>
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className={cc.form.label}>
                Nieuw Wachtwoord <span className="text-red-500">*</span>
              </label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={cn(
                  'w-full px-4 py-3 rounded-xl border-2',
                  'focus:outline-none focus:ring-2 focus:ring-primary/20',
                  'border-gray-200 focus:border-primary'
                )}
                placeholder="Minimaal 8 karakters"
                required
                minLength={8}
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className={cc.form.label}>
                Bevestig Wachtwoord <span className="text-red-500">*</span>
              </label>
              <input
                type="password"
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className={cn(
                  'w-full px-4 py-3 rounded-xl border-2',
                  'focus:outline-none focus:ring-2 focus:ring-primary/20',
                  'border-gray-200 focus:border-primary'
                )}
                placeholder="Herhaal je wachtwoord"
                required
                minLength={8}
              />
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-sm text-red-800">{error}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={isSubmitting}
              className={cn(
                'w-full px-6 py-3 font-semibold text-white rounded-lg',
                colors.primary.bg,
                colors.primary.hover,
                'disabled:opacity-50 disabled:cursor-not-allowed',
                cc.transition.base
              )}
            >
              {isSubmitting ? 'Bezig met upgraden...' : 'Upgrade naar Volledig Account'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Al een volledig account?{' '}
              <a href="/login" className={cn(colors.primary.text, 'underline')}>
                Log hier in
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
```

---

### STAP 7: API Endpoints Configuratie

**Bestand:** `src/services/api/endpoints.ts`

Voeg toe:
```typescript
export const API_ENDPOINTS = {
  // Bestaande endpoints...
  
  // 🆕 NIEUW: Publieke registratie endpoints (V30)
  publicRegistration: '/api/public/aanmelden',
  upgradeAccount: '/api/public/upgrade-to-full-account',
  activeEvent: '/api/public/events/active',
  
  // Oude endpoints (voor backwards compatibility)
  participants: '/api/participant', // Nu alleen voor admin
  eventRegistrations: '/api/event-registrations', // Nu alleen voor admin
};
```

---

## ✅ Checklist Frontend Implementatie

### Fase 1: Basis Wijzigingen
- [ ] Update `schema.ts` met nieuwe velden en validaties
- [ ] Maak `AccountTypeSelector.tsx` component
- [ ] Maak `PasswordField.tsx` component
- [ ] Update imports in `FormContainer.tsx`

### Fase 2: FormContainer Aanpassingen
- [ ] Voeg nieuwe state toe (`accountChoice`, `showPasswordField`)
- [ ] Voeg `AccountTypeSelector` toe in form flow
- [ ] Voeg conditioneel `PasswordField` toe
- [ ] Update `onSubmit` om nieuwe endpoint te gebruiken
- [ ] Update error handling voor nieuwe error codes
- [ ] Test alle validaties

### Fase 3: Success Message Aanpassingen
- [ ] Update `SuccessMessageProps` interface
- [ ] Voeg conditionele content toe per account type
- [ ] Voeg App Download sectie toe (alleen full account)
- [ ] Voeg Upgrade CTA toe (alleen temporary)

### Fase 4: Upgrade Flow
- [ ] Maak `UpgradeAccount.tsx` pagina
- [ ] Voeg route toe in React Router (`/upgrade`)
- [ ] Implementeer upgrade form
- [ ] Implementeer error handling
- [ ] Test upgrade flow end-to-end

### Fase 5: UI/UX Polish
- [ ] Voeg animaties toe bij account type switch
- [ ] Voeg password strength indicator toe
- [ ] Voeg tooltips/help text toe waar nodig
- [ ] Test accessibility (keyboard navigation, screen readers)
- [ ] Test responsive design (mobile, tablet, desktop)

### Fase 6: Testing
- [ ] Test full account registratie flow
- [ ] Test temporary registratie flow
- [ ] Test upgrade flow
- [ ] Test alle validatie errors
- [ ] Test duplicate registratie scenarios
- [ ] Test met verschillende rollen en afstanden

---

## 🎨 UI/UX Aanbevelingen

### Account Type Keuze
- **Positie:** Direct NA contactgegevens, VOOR rol selectie
- **Reden:** Gebruiker moet eerst kiezen of ze een account willen voordat ze verder gaan
- **Visual:** Gebruik cards met duidelijke voor/nadelen

### Wachtwoord Veld
- **Toon alleen als:** `want_account === true`
- **Animatie:** Smooth fade-in bij tonen, fade-out bij verbergen
- **Positie:** Direct NA account type keuze
- **Features:** 
  - Toggle visibility (oog icoon)
  - Strength indicator (rood/geel/groen)
  - Helper text met tips

### Validatie Feedback
- **Real-time validation:** Bij onBlur of onChange
- **Error messages:** Duidelijk en actionable
- **Success indicators:** Groene checkmarks bij valide velden

### Upgrade CTA
- **Temporary registratie success:** Prominente CTA om te upgraden
- **Email:** Link naar upgrade pagina in confirmation email
- **Timing:** Niet te pushy, maar wel duidelijk zichtbaar

---

## 🔄 API Error Handling

### Error Codes en Messages

| Code | HTTP Status | Gebruiker Bericht | Actie |
|------|-------------|-------------------|-------|
| `ALREADY_REGISTERED` | 409 | "Je bent al inges chronven voor dit jaar" | Toon login link |
| `EMAIL_EXISTS` | 409 | "Er bestaat al een account met dit e-mailadres" | Toon login/reset link |
| `VALIDATION_ERROR` | 400 | Specifieke validatie fout | Highlight veld |
| `NO_TEMPORARY_ACCOUNT` | 404 | "Geen tijdelijke registratie gevonden" | Toon registratie link |
| `ALREADY_FULL_ACCOUNT` | 400 | "Dit is al een volledig account" | Toon login link |
| `WEAK_PASSWORD` | 400 | "Wachtwoord te zwak" | Highlight wachtwoord |
| `REGISTRATION_FAILED` | 500 | "Er ging iets mis, probeer later opnieuw" | Retry button |

**Implementatie:**
```typescript
const handleApiError = (error: any): string => {
  const code = error?.response?.data?.code;
  const message = error?.response?.data?.error;
  
  const errorMap: Record<string, string> = {
    'ALREADY_REGISTERED': 'Je bent al ingeschreven voor dit jaar. Wil je je gegevens wijzigen? Neem contact op.',
    'EMAIL_EXISTS': 'Er bestaat al een account met dit e-mailadres. Probeer in te loggen of gebruik een ander e-mailadres.',
    'NO_TEMPORARY_ACCOUNT': 'Geen tijdelijke registratie gevonden. Heb je je al geregistreerd?',
    'ALREADY_FULL_ACCOUNT': 'Dit account is al een volledig account! Je kunt inloggen in de app.',
    'WEAK_PASSWORD': 'Wachtwoord moet minimaal 8 karakters bevatten.',
    'VALIDATION_ERROR': message || 'De ingevoerde gegevens zijn ongeldig.',
  };
  
  return errorMap[code] || message || 'Er ging iets mis. Probeer het later opnieuw.';
};
```

---

## 📱 Volgorde van Formulier Secties

**Aanbevolen volgorde voor optimal UX:**

1. **Contactgegevens** (Naam, Email)
2. **🆕 Account Type Keuze** (Full vs Temporary)
3. **🆕 Wachtwoord** (Conditioneel - alleen bij Full)
4. **Rol Selectie** (Deelnemer, Begeleider, Vrijwilliger)
5. **Telefoon** (Conditioneel - alleen bij Begeleider/Vrijwilliger)
6. **Afstand Selectie** (2.5km - 15km)
7. **Ondersteuning** (Ja, Nee, Anders)
8. **Bijzonderheden** (Conditioneel - bij Ondersteuning Ja/Anders)
9. **Algemene Voorwaarden**
10. **Submit Button**

**Rationale:** Account type keuze vroeg in flow zodat gebruiker weet wat te verwachten (wachtwoord veld wel/niet).

---

## 🧪 Frontend Testplan

### Unit Tests
```typescript
describe('AccountTypeSelector', () => {
  it('renders both options', () => {});
  it('calls onChange with correct value', () => {});
  it('shows error message if present', () => {});
  it('highlights selected option', () => {});
});

describe('PasswordField', () => {
  it('shows/hides password on toggle', () => {});
  it('calculates password strength correctly', () => {});
  it('shows validation errors', () => {});
});

describe('FormContainer registration flow', () => {
  it('shows password field when full account selected', () => {});
  it('hides password field when temporary selected', () => {});
  it('validates password when full account selected', () => {});
  it('does not require password for temporary', () => {});
  it('submits correct data to API', () => {});
});
```

### Integration Tests
```typescript
describe('Full Account Registration E2E', () => {
  it('completes full registration with all fields', async () => {
    // Fill form
    // Select full account
    // Enter password
    // Submit
    // Expect success with gebruiker_id
  });
});

describe('Temporary Registration E2E', () => {
  it('completes temporary registration', async () => {
    // Fill form
    // Select temporary
    // Submit (no password)
    // Expect success without gebruiker_id
  });
});

describe('Account Upgrade E2E', () => {
  it('upgrades temporary to full account', async () => {
    // Register temporary
    // Navigate to upgrade
    // Enter email + password
    // Submit
    // Expect success with heeft_id
  });
});
```

---

## 📞 Support & Troubleshooting

### Veelvoorkomende Issues

**Issue:** "Wachtwoord veld wordt niet getoond"  
**Fix:** Check of `want_account` state correct is gezet. Debug met React DevTools.

**Issue:** "Form valideert niet correct"  
**Fix:** Check Zod schema custom refine() functies. Log validatie errors.

**Issue:** "API geeft 400 error"  
**Fix:** Check of alle required velden worden meegestuurd. Compare met backend validation.

**Issue:** "Email wordt niet verzonden"  
**Fix:** Check backend logs. Email service moet draaien.

### Debug Logging

Voeg toe aan `onSubmit`:
```typescript
console.log('Submitting registration:', {
  ...validatedData,
  wachtwoord: validatedData.wachtwoord ? '***' : undefined // Never log actual password
});
```

---

## 🎬 Demo Flows

### Demo 1: Full Account (Happy Path)
1. Open `/aanmelden`
2. Vul naam in: "Demo Gebruiker"
3. Vul email in: "demo@example.com"
4. Kies: "Ja, maak een account aan"
5. **Wachtwoord veld verschijnt**
6. Vul wachtwoord in: "DemoPass2026"
7. Kies rol: "Deelnemer"
8. Kies afstand: "10 KM"
9. Kies ondersteuning: "Nee"
10. Accept terms
11. Submit
12. **Success met app download links**

### Demo 2: Temporary (Happy Path)
1. Open `/aanmelden`
2. Vul naam in: "Quick Register"
3. Vul email in: "quick@example.com"
4. Kies: "Nee, alleen registreren"
5. **Geen wachtwoord veld**
6. Kies rol: "Deelnemer"
7. Kies afstand: "6 KM"
8. Kies ondersteuning: "Nee"
9. Accept terms
10. Submit
11. **Success met upgrade CTA**

### Demo 3: Upgrade Flow
1. Open `/upgrade`
2. Vul email in: "quick@example.com"
3. Vul wachtwoord in: "NewPass2026"
4. Bevestig wachtwoord: "NewPass2026"
5. Submit
6. **Success toast + redirect naar app download**

---

## 📊 Analytics Events

Track deze events voor insights:

```typescript
// Account type selectie
logEvent('registration', 'select_account_type', 'full' | 'temporary');

// Password field interaction
logEvent('registration', 'password_field_interaction', 'show' | 'hide');

// Submission per type
logEvent('registration', 'submit_full_account', `${rol}_${afstand}`);
logEvent('registration', 'submit_temporary', `${rol}_${afstand}`);

// Upgrade
logEvent('upgrade', 'upgrade_attempt', email);
logEvent('upgrade', 'upgrade_success', email);
logEvent('upgrade', 'upgrade_fail', error_code);

// CTA clicks
logEvent('registration', 'click_upgrade_cta', 'success_page');
```

---

## 🚦 Deployment Stappen

### Pre-Deployment
1. ✅ Backend V30 gedeployed en getest
2. ✅ Database migratie uitgevoerd
3. ✅ Email templates beschikbaar
4. ✅ API endpoints werkend getest

### Frontend Deployment
1. Implementeer alle wijzigingen in development branch
2. Test lokaal met backend (docker-compose up)
3. Test alle flows (full, temporary, upgrade)
4. Code review
5. Merge naar main
6. Deploy naar staging
7. QA testen op staging
8. Deploy naar production

### Post-Deployment
1. Monitor error rates
2. Check registratie metrics
3. Verify emails worden verstuurd
4. Test app login voor nieuwe full accounts
5. Monitor upgrade conversie rate

---

## 📧 Email Integratie

Backend stuurt automatisch emails via templates:

1. **Full Account:** `registration_full_account_email.html`
   - Welkomstbericht
   - App download links  
   - Inloggegevens

2. **Temporary:** `registration_temporary_account_email.html`
   - Bevestiging
   - Event info
   - **Upgrade CTA link**

3. **Upgrade:** `account_upgrade_email.html`
   - Upgrade bevestiging
   - App download links
   - Nieuwe mogelijkheden

Frontend hoeft **GEEN** emails te verzenden - backend handled dit!

---

## 🔗 Backend API Reference

### Endpoint: POST /api/public/aanmelden

**Request:**
```json
{
  "naam": "string",
  "email": "string (email format)",
  "telefoon": "string | undefined",
  "rol": "Deelnemer" | "Begeleider" | "Vrijwilliger",
  "afstand": "2.5 KM" | "6 KM" | "10 KM" | "15 KM",
  "ondersteuning": "Ja" | "Nee" | "Anders",
  "bijzonderheden": "string",
  "want_account": boolean,
  "wachtwoord": "string | undefined (required if want_account=true, min 8 chars)",
  "terms": boolean (must be true),
  "test_mode": boolean
}
```

**Response 201:**
```json
{
  "success": true,
  "message": "String met bevestiging",
  "participant_id": "uuid",
  "registration_id": "uuid",
  "account_type": "full" | "temporary",
  "has_app_access": boolean,
  "gebruiker_id": "uuid (optional - only for full accounts)",
  "event_name": "De Koninklijke Loop 2026",
  "event_date": "16-05-2026"
}
```

### Endpoint: POST /api/public/upgrade-to-full-account

**Request:**
```json
{
  "email": "string (must match existing temporary account)",
  "wachtwoord": "string (min 8 chars)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "Je account is geüpgraded!",
  "gebruiker_id": "uuid",
  "participant_id": "uuid",
  "has_app_access": true
}
```

---

## 💡 Best Practices

### State Management
- Gebruik `useState` voor UI state (`showPasswordField`)
- Gebruik `react-hook-form` voor form state
- Sync password field visibility met `useEffect`

### Performance
- Gebruik `memo()` voor components die niet vaak re-renderen
- Lazy load upgrade pagina (React.lazy + Suspense)
- Debounce password strength calculation

### Accessibility
- Alle interactive elements keyboard accessible
- Screen reader friendly labels
- Error messages hebben `role="alert"`
- Focus management bij modal/section changes

### Error Handling
- Friendly error messages (geen technische jargon)
- Suggesties voor oplossing (bijv. "Probeer in te loggen" bij EMAIL_EXISTS)
- Retry mogelijkheid bij failures

---

## 🎓 Training Material

### Voor Support Team

**Veelgestelde Vragen:**

Q: Wat is het verschil tussen de account types?  
A: Full account = app toegang + permanent. Temporary = alleen evenement 2026, geen app.

Q: Kan iemand later upgraden?  
A: Ja! Via /upgrade pagina of link in confirmation email.

Q: Wat als iemand zijn wachtwoord vergeet?  
A: Gebruik de "Wachtwoord vergeten" link op de login pagina.

Q: Kunnen temporary accounts inloggen in de app?  
A: Nee, alleen full accounts hebben app toegang.

---

## 📞 Contact

Bij vragen over de implementatie:
- **Backend:** Zie `docs/V30_DUAL_REGISTRATION_SYSTEM.md`
- **API Reference:** Zie `docs/api/` folder
- **Issues:** Meld via Linear of Github Issues

---

**Last Updated:** 2025-11-10  
**Next Review:** Bij eerste deployment naar production