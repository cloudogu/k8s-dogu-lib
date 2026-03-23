# Dogu format

Die Dogu-CR kann genutzt werden, um Cloudogu-Dogus in einem Kubernetescluster mit dem Dogu-Operator zu installieren.
Es können verschiedene Einstellungen getroffen werden, die zusätzlich zur dogu.json und der Dogu-Konfiguration
ausgewertet werden.

Folgend werden alle Felder einer Dogu-CR beschrieben und mit Beispielen veranschaulicht.

## Komplettes Beispiel

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Dogu
metadata:
  name: postfix
spec:
  name: official/usermgt
  version: 1.20.0-5
  additionalIngressAnnotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
  additionalMounts:
    - sourceType: ConfigMap
      name: my-configmap
      volume: importHistory
      subfolder: "my-configmap-subfolder"
  exportMode: false
  resources:
    minDataVolumeSize: 2Gi
  security:
    appArmorProfile:
      localhostProfile: "localhost-profile"
      type: Localhost
    capabilities:
      add:
        - CAP_AUDIT_CONTROL
      drop:
        - CAP_SETGID
    readOnlyRootFileSystem: false
    runAsNonRoot: false
    seLinuxOptions:
      level: internal
      role: user_r
      type: user_t
      user: user_u
    seccompProfile:
      type: RuntimeDefault
  stopped: false
  supportMode: false
  upgradeConfig:
    allowNamespaceSwitch: false
    forceUpgrade: false
```

## Name

* Pflichtfeld
* Datentyp: string
* Inhalt: Gibt den Namen einschließlich des Namespace des Dogu an.
* Beispiel: `"name": "official/usermgt"`

## Version

* Pflichtfeld
* Datentyp: string
* Inhalt: Gibt die Version des Dogu an.
* Beispiel: `"version": "1.20.0-5"`

## AdditionalIngressAnnotations

* Optional
* Datentyp: string
* Inhalt: AdditionalIngressAnnotations liefert zusätzliche Anmerkungen, die in die Ingress-Regeln des Dogus aufgenommen
  werden.
* Beispiel:

```
additionalIngressAnnotations:
  nginx.ingress.kubernetes.io/proxy-body-size: "0"
```

## AdditionalMounts

* Optional
* Datentyp: Array<DataMount>
* Inhalt: AdditionalMounts bietet die Möglichkeit, zusätzliche Daten in das Dogu einzubinden.
* Beispiel:

```
  additionalMounts:
    - sourceType: ConfigMap
      name: my-configmap
      volume: importHistory
      subfolder: "my-configmap-subfolder"
    - sourceType: Secret
      name: my-secret
      volume: importHistory
```

### DataMount

Ein DataMount kann die folgenden Felder enthalten:

#### SourceType

* Pflichtfeld
* Datentyp: Enum <ConfigMap; Secret>
* Inhalt: SourceType legt fest, woher die Daten stammen.
  Gültige Optionen sind:
    - ConfigMap - Daten, die in einer kubernetes ConfigMap gespeichert sind.
    - Secret - Daten, die in einem kubernetes Secret gespeichert sind.
* Beispiel: `"sourceType": ConfigMap`

#### Name

* Pflichtfeld
* Datentyp: String
* Inhalt: Name ist der Name der Datenquelle.
* Beispiel: `"name": my-configmap`

#### Volume

* Pflichtfeld
* Datentyp: String
* Inhalt: Volume ist der Name des Volumes, in das die Daten gemountet werden sollen. Dieses wird in der jeweiligen
  dogu.json definiert.
* Beispiel: `"volume": importHistory`

#### Subfolder

* Optional
* Datentyp: String
* Inhalt: Subfolder definiert einen Unterordner, in dem die Daten innerhalb des Volumes abgelegt werden sollen.
* Beispiel: `"subfolder": "my-configmap-subfolder"`

## ExportMode

* Optional
* Datentyp: String
* Inhalt: ExportMode gibt an, ob sich das Dogu im „Exportmodus“ befinden soll. Wenn dies der Fall ist, erzeugt der
  Operator einen Sidecar-Container zusammen mit einem neuen Volume-Mount, um den Migrationsprozess von einem Cloudogu
  EcoSystem zu einem anderen zu unterstützen.
* Beispiel: `"exportMode": false`

## Resources

* Optional
* Datentyp: Object
* Inhalt: Ressourcen des Dogus (e.g. minDataVolumeSize)
* Beispiel:

```
resources:
  minDataVolumeSize: 2Gi
```

### MinDataVolumeSize

* Optional
* Datentyp: String
* Inhalt: MinDataVolumeSize stellt die gewünschte Mindestgröße des Volumes dar. Das Erhöhen dieses Wertes führt ggf. zu einer
  automatischen Erweiterung. Dies beinhaltet eine Ausfallzeit für die jeweilige Dogu. Die Standardgröße für Volumes ist
  „2Gi“. Es ist nicht möglich, die Größe des Volumes nach einer Erweiterung zu verringern. Dies würde zu einem
  inkonsistenten Zustand des Dogus führen.
* Beispiel: `"minDataVolumeSize": 2Gi`

## Security

* Optional
* Datentyp: Object
* Inhalt: Security überschreibt die im Dogu-Deskriptor definierten Sicherheitsrichtlinien. Diese Felder können verwendet
  werden, um die Angriffsfläche einer Dogu weiter zu reduzieren.

Security kann die folgenden Attribute enthalten:

### AppArmorProfile

* Optional
* Datentyp: Object
* Inhalt: AppArmorProfile sind die von diesem Container zu verwendenden AppArmor-Optionen.
* Beispiel:

```
appArmorProfile:
  localhostProfile: "localhost-profile"
  type: Localhost
```

#### Type

* Pflichtfeld
* Datentyp: Enum <Localhost; RuntimeDefault; Unconfined>
* Inhalt: Type gibt an, welche Art von AppArmor-Profil angewendet wird.
  Gültige Optionen sind:
    - Localhost - ein auf dem Node vorgeladenes Profil.
    - RuntimeDefault - das Standardprofil der Container-Runtime.
    - Unconfined - keine AppArmor-Durchsetzung.

#### LocalhostProfile

* Optional
* Datentyp: String
* Inhalt: LocalhostProfile gibt ein auf dem Node geladenes Profil an, das verwendet werden soll.
  Das Profil muss auf dem Node vorkonfiguriert sein, um zu funktionieren. Der Name muss mit dem des geladenen Profils
  übereinstimmen.
  Muss nur gesetzt werden, wenn der Typ „Localhost“ ist.

### Capabilities

* Optional
* Datentyp: Object
* Inhalt: Capabilities legt die erlaubten und nicht erlaubten Fähigkeiten für das Dogu fest. Das Dogu sollte nicht mehr
  als die hier konfigurierten Capabilities verwenden, da es sonst beim Start oder zur Laufzeit zu Fehlern kommen kann.
  Jede Fähigkeit repräsentiert einen POSIX-„Capabilities“-Typ. Siehe Dokumentationen
  unter https://manned.org/capabilities.7
* Beispiel:

```
capabilities:
  add:
    - CAP_AUDIT_CONTROL
  drop:
    - CAP_SETGID
```

#### Add

* Optional
* Datentyp: Array<String>
* Inhalt: Add enthält die Capabilities, deren Verwendung in einem Container erlaubt sein soll. Diese Liste ist optional.

#### Drop

* Optional
* Datentyp: Array<String>
* Inhalt: Drop enthält die Capabilities, die für die Verwendung in einem Container gesperrt werden sollen. Diese Liste
  ist optional.

### ReadOnlyRootFileSystem

* Optional
* Datentyp: boolean
* Inhalt: ReadOnlyRootFileSystem mountet das Root-Dateisystem des Containers als schreibgeschützt. Das Dogu muss den
  Zugriff auf das Root-Dateisystem nur lesend unterstützen, ansonsten kann der Dogu-Start fehlschlagen. Dieses Flag ist
  optional und steht standardmäßig auf nil.
  Ist es gleich nil, wird der im Dogu-Deskriptor definierte Wert verwendet.
* Beispiel: `"readOnlyRootFileSystem": true`

### RunAsNonRoot

* Optional
* Datentyp: boolean
* Inhalt: RunAsNonRoot gibt an, dass der Container als Nicht-Root-User laufen muss. Das Dogu muss die Ausführung als
  Nicht-Root-User unterstützen, andernfalls kann der Start vom Dogu fehlschlagen. Dieses Flag ist optional und steht
  standardmäßig auf nil.
  Ist es gleich nil, wird der im Dogu-Deskriptor definierte Wert verwendet.
* Beispiel: `"runsNonRoot": true`

### SeLinuxOptions

* Optional
* Datentyp: Object
* Inhalt: SELinuxOptions ist der SELinux-Kontext, der auf den Container angewendet werden soll.
  Wenn nicht angegeben, wird die Container-Runtime einen zufälligen SELinux- Kontext für jeden Container zuweisen, was
  das kubernetes Standardverhalten ist.
* Beispiel:

```
seLinuxOptions:
  level: internal
  role: user_r
  type: user_t
  user: user_u
```

#### Level

* Optional
* Datentyp: string
* Inhalt: Level ist das SELinux-Level-Label, das für den Container gilt.

#### Role

* Optional
* Datentyp: string
* Inhalt: Role ist das SELinux-Role-Label, das für den Container gilt.

#### Type

* Optional
* Datentyp: string
* Inhalt: Type ist das SELinux-Type-Label, das für den Container gilt.

#### User

* Optional
* Datentyp: string
* Inhalt: User ist das SELinux-User-Label, das für den Container gilt.

### SeccompProfile

* Optional
* Datentyp: Object
* Inhalt: SeccompProfile sind die Seccomp-Optionen, die von diesem Container verwendet werden sollen.
* Beispiel:

```
seccompProfile:
  localhostProfile: "localhost-profile"
  type: Localhost
```

#### Type

* Pflichtfeld
* Datentyp: Enum <Localhost; RuntimeDefault; Unconfined>
* Inhalt: Typ gibt an, welche Art von SeccompProfile angewendet wird.
  Gültige Optionen sind:
    - Localhost - es soll ein Profil verwendet werden, das in einer Datei auf dem Node definiert ist.
    - RuntimeDefault - es soll das Standardprofil der Container- Runtime verwendet werden.
    - Unconfined - es soll kein Profil angewendet werden.

#### LocalhostProfile

* Optional
* Datentyp: string
* Inhalt: LocalhostProfile gibt an, dass ein in einer Datei auf dem Node definiertes Profil verwendet werden soll.
  Das Profil muss auf dem Node vorkonfiguriert sein, damit es funktioniert.
  Es muss ein Pfad sein, der relativ zum Speicherort des konfigurierten seccomp-Profils des Kubelet ist.
  Muss gesetzt werden, wenn der Typ „Localhost“ ist. Darf NICHT für einen anderen Typ gesetzt werden.

## Stopped

* Optional
* Datentyp: boolean
* Inhalt: Stopped gibt an, ob das Dogu laufen soll (stopped=false) oder nicht (stopped=true).
* Beispiel: `"stopped": true`

## SupportMode

* Optional
* Datentyp: boolean
* Inhalt: SupportMode gibt an, ob das Dogu im Support-Modus neu gestartet werden soll (z. B. um manuell aus eine
  Absturzschleife zu beheben).
* Beispiel: `"supportMode": true`

## UpgradeConfig

* Optional
* Datentyp: Object
* Inhalt: UpgradeConfig enthält Optionen zur Beeinflussung des Upgrade-Prozesses.
* Beispiel:

```
upgradeConfig:
  allowNamespaceSwitch: false
  forceUpgrade: false
```

### AllowNamespaceSwitch

* Optional
* Datentyp: boolean
* Inhalt: AllowNamespaceSwitch lässt ein Dogu seinen Dogu-Namespace während eines Upgrades wechseln. Das Dogu muss
  technisch dasselbe Dogu sein, das sich in einem anderen Namespace befunden hat. Die Version des entfernten Dogus muss
  gleich oder größer als die Version des lokalen Dogus sein.

### ForceUpgrade

* Optional
* Datentyp: boolean
* Inhalt: ForceUpgrade erlaubt es, die gleiche oder sogar eine niedrigere Dogu-Version zu installieren, als bereits
  installiert ist. Bitte beachten Sie, dass durch ein unsachgemäßes Dogu-Downgrade Datenverluste auftreten können.
