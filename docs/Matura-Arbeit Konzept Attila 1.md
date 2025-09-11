# Titel: Wirtschaftssimulation als Videospiel
## Einleitung
Für meine Maturaarbeit möchte ich ein Videospiel programmieren, in dem man als Geschäftsleitung einer Firma zusammen oder gegen Andere spielen kann.

In diesem Spiel werden Firmen ein Produkt herstellen und verkaufen müssen und gegeneinander konkurrieren um die grösste Firma* zu werden. In der Simulation werden verschiedene Kaufverhalten und (zum Teil) äussere Umweltfaktoren simuliert. So werden die Firmen sich anpassen müssen um ihre Konkurrenz zu besiegen.

­\* Die Gewinnkonditionen sind noch nicht festgelegt, es könnte mehrere geben.
## Motivation
Während der Wirtschaftswoche haben wir mit einer Simulation gearbeitet, damit wir lernen, wie man einen Unternehmen führt und welche Challenges ein Unternehmen haben könnte. Das hat mich komplett fasziniert: nicht nur die Simulation sondern auch das Spielen mit Kollegen und die Konkurrenz. Ich habe mir direkt vorgestellt sowas mit Freunde zu spielen, in mein Freizeit.
Da ich auch Ziele in Finanz und Softwareengineering habe, scheint das die perfekte Kombinationen von Fächern zu sein.

# Leitfrage

Kann man innerhalb der vorgegebenen Zeit eine online/multiplayer betriebswirtschaftliche Simulation bauen, die sowohl eine Herausforderung für den Spieler ist und Spass macht?
# Erarbeitung

Das Spiel wird sich um ein Wirtschaftliches Modell (Siehe [[Modell.canvas|Modell]]) drehen. Das Modell wird erst grob gebaut, dann optimiert, um Spielspass und Realismus zu maximieren. Da mich das Projekt eher auf der Programmierseite interessiert, werde ich nicht zu genau auf die wirtschaftliche Recherche eingehen.

Das Spiel wird aus 3 verschieden Prozessen (Programme) bestehen:
- Der Server: Im Server wird die Simulation laufen, darin werden alle (primären**) Rechnungen durchgeführt.
- Der Client: Ein Interface, damit die einzelnen Spieler die Informationen vom Server bekommen können.
- Das GUI / User Interface: Eine GUI (Graphical User Interface), damit ein Spieler mit der Simulation interagieren kann.

\*\* Damit meine ich Simulations-relevante Rechnungen, es können Rechnungen auf dem Client durchgeführt werden, die sich aus den Simulationswerte ableiten.

Um die Entwicklungszeit zu reduzieren werde ich in `Go` programmieren auf Grund meiner Erfahrung in C, das vereinfachte Memory-Management und der hohen Performance.
# Produkte der Arbeit
- Ein Entwicklungstagebuch
- Eine betriebswirtschaftliche Simulation/Spiel
# Zeitplan
## 1. Modell aufbauen ([[Modell.canvas|Modell]])
Ich notiere, wie das Modell aussieht
## 2. Simulation mit Modell implementieren
Während den Frühlingssemester werde ich die Simulation implementieren
## 3. Client & GUI bauen
In der 1. Hälfte des Herbstsemesters mache ich die GUI fertig.
## 4 Networking bauen
## 5. Testen, Balancing & Abschliessen

## 6. Dokumentation & Präsentieren