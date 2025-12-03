# Wieso eine Spezifische Funktion?
In der Simulation werden alle mögliche Kunden ein Einkommen zugeteilt. Das Einkommen hat einen Einfluss auf wie viel die Kunde zahlen kann.

Damit die Simulation realistischer ist, wollte ich eine realistische Verteilung der einkommen.
# Wie wird das gemacht
Wenn man die echte Verteilung von einkommen anschaut, ähnelt das eine Normalverteilung, die verschoben ist.
![[Normalverteilung Diagramm#^area=hRYNKNhIX7-FJjHSUCjGH|center]]
Am Anfang wollte ich direkt Proben von Normalverteilungen generieren, z.B. mit der Box-Muller Transformation ([wikipedia](https://en.wikipedia.org/wiki/Box%E2%80%93Muller_transform)). Durch proben habe ich festgestellt, dass die Box-Muller Transformation zu inakkurat und kompliziert war. Ich fand später raus, dass es in Go schon eine Funktion gab, die eine Normalprobe generieren kann: `NormFloat64()`.

```Go
func randincome(mean_income int, standard_dev int) int {
	income := -1
	for income > 1000 {
		income = int(rand.NormFloat64()*standard_dev) + mean_income
	}
	return income
}
```

Die Funktion sieht am Schluss so aus. Es tut einfach eine Normalverteilung abschneiden falls der Wert zu tief ist.
![](../../attachments/Pasted%20image%2020250820173943.png)

*Verteilung von Einkommen nach der Funktion*