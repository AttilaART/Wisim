# Titel: Wirtschaftssimulation als Videospiel
## Einleitung
Für meine Maturaarbeit möchte ich ein Videospiel programmieren, indem man als Geschäftsleitung einer Firma, zusammen oder gegen anderen, spielen kann.
In diesem Spiel werden Firmen einen Produkt herstellen und verkaufen müssen und gegen einander konkurrieren um die Grösste Firma* zu werden. In der Simulation werden verschiedene Kaufverhalten und (zum teil) äussere Umweltfaktoren simuliert. So werden die Firmen sich anpassen müssen um ihre Konkurrenz zu besiegen.

­\* Die Gewinnkonditionen sind noch nicht festgelegt, es könnte mehrere geben.
## Motivation
Während der Wirtschaftswoche haben wir mit einer Simulation gearbeitet, damit wir lernen, wie man einen Unternehmen führt und welche Challenges ein Unternehmen haben könnte. Das hat mich komplett fasziniert: nicht nur die Simulation sondern auch das Spielen mit Kollegen und die Konkurrenz. Ich habe mir direkt vorgestellt sowas mit Freunde zu spielen in mein Freizeit.
Da ich auch Ziele in Finanz und Softwareengineering habe, scheint das die perfekte Kombinationen von Fächern zu sein.

# Leitfrage

Kann man eine online/multiplayer betriebswirtschaftliche Simulation bauen, die sowohl Spass macht zu spielen und einigermassen akkurat eine Wirtschaft simuliert innerhalb der vorgegebenen Zeit?
# Erarbeitung
## Theorie

Um eine sinnvolle Simulation zu bauen, muss ich folgende Faktoren untersuchen:
- Kaufverhalten
	- und den Einfluss von Marketing darauf
	- in verschiedene Einkommensklassen
	- den Einfluss von äussere Umweltfaktoren
	- den Einfluss von der ökologische Wahrnehmung eines Produktes auf Verhalten.
- Finanzen
	- wie Bilanzen und Erfolgsrechnungen (unter anderem) berechnet werden.
	- Wie Budgets berechnet werden.
	- Was die wichtigsten Werte sind für die Unternehmensleitung
	- Welche Faktoren angeschaut werden, wenn investiert wird
	- Wie Banken Kredite ausgeben.
- Mitarbeiter
	- Produktivität von Arbeiter nach:
		- Kultur
		- Lohn
		- etc.
	- Effekte von Berufsicherheit / -unsicherheit.
	- Wie / wenn / warum Mitarbeiter Firma wechseln.

Für den spielerischen Aspekt muss ich vor allem Networking und Modellierung anschauen. Z. B. muss ich schauen, ob ich ein statistischer Modell mache, wo Werte einfach statistisch ausgerechnet werden, oder ob ich einen Agent-basierter Modell mache, wo die Simulation aus viele unabhängige Personen ("Agenten"), die alle selbst entscheiden können, besteht.
## Entwickelung
Wenn ich mein Modell habe, muss ich 3 verschieden Prozesse (Programme) entwickeln:
- Der Server: Im Server wird die Simulation laufen, drin werden alle (primäre**) Rechnungen durchgeführt.
- Der Client: Ein Interface, damit die einzelnen Spieler die Informationen vom Server bekommen können.
- Das GUI / User Interface: Einen GUI (Graphical User Interface), damit ein Spieler mit der Simulation interagieren kann.

\*\* Damit meine ich Simulation-relevante Rechnungen, es können Rechnungen auf dem Client durchgeführt werden, die sich aus Simulationswerte ableiten.

# Produkte der Arbeit
- Ein Entwickelungstagebuch
- Eine Auseinandersetzung der Oben (unter Theorie) genannte Faktoren
- Eine betriebswirtschaftliche Simulation
# Zeitplan
## 1. Recherche-Phase
Während meiner Freizeit im Frühlingssemester 25 werde ich vor allem Bücher, die sich mit dem Thema befassen, auseinandersetzen.
## 2. Modell aufbauen ([[Modell.canvas|Modell]])
Teilweise während Phase 1 werde ich notieren, wie das Modell aussehen soll.
## 3. Simulation mit Modell implementieren
Während den Sommerferien und den Anfang vom Herbstsemester werde ich die Simulation implementieren
## 4 Networking bauen
Damit der Client überhaupt funktioniert, bereite ich schon im Voraus die Grundlagen vom Networking vor
## 5. Client & GUI bauen
In der 2. Hälfte des Herbstsemesters mache ich das GUI fertig.
## 6. Testen & Abschliessen

## 7. Dokumentation & Präsentieren