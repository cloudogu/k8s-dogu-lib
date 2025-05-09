# Dogurestart format

Die Dogurestart-CR kann genutzt werden, um Dogus explizit neustarten zu können. Der Dogu-Operator skaliert daraufhin
dsa entsprechende Dogu-Deployment herunter und anschließend wieder hoch. 

Folgend werden alle Felder einer Dogurestart-CR beschrieben und mit Beispielen veranschaulicht.

## Komplettes Beispiel

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: DoguRestart
metadata:
  generateName: usermgt-restart-
spec:
 doguName: usermgt
```

Bitte beachten: `generateName` kann genutzt werden, um einen eindeutigen Namen zu erzeugen. Dies funktioniert jedoch 
nicht mit `kubectl apply` sondern nur `kubectl create`

## doguName

* Pflichtfeld
* Datentyp: string
* Inhalt: Das Feld `doguName` gibt den Namen des Dogus an, welches neugestartet werden soll.
* Beispiel: `"doguName": "usermgt"`
