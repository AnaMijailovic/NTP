# Analiza i organizacija fajlova i foldera

### Autor
    Ana Mijailović, SW13/2016
    mijailovic.sw13.2016@uns.ac.rs

### 1. Analiza i vizualizacija fajlova i foldera
      a) Prikaz u vidu stabla počevši od odabrane putanje
      b) Bar Chart i Pie Chart prikazi zauzeća memorije po:
          - Tipovima fajla
          - Datumu kreiranja (po godinama ili mesecima)
          - Datumu poslednje izmene (po godinama ili mesecima)
          - Datumu poslednjeg pristupa fajlu (po godinama ili mesecima)

### 2. Reorganizacija strukture foldera
       Mogućnost organizovanja fajlova na zadatoj putanji (rekurzivno ili ne) tako što bi se odvojili u 
       posebne foldere na osnovu:
         a) Tipa fajla
         b) Veličine fajla (korisnik zadaje korak)
         c) Datuma kreiranja (pun datum, mesec i godina ili samo godina)
         
         *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje 
            originalne strukture
            
### 3. Preimenovanje svih foldera/fajlova na zadatoj putanji
        a) Generisanjem random naziva
        b) Zadavanjem prefiksa/sufiksa koji se dodaju na postojeći naziv
        c) Uklanjanjem zadatog dela naziva ili zamenom tog dela naziva novim
        d) Zadavanjem izraza gde će delovi navedeni unutar vitičastih zagrada
           biti zamenjeni odgovarajućim sadržajem:
                -	{name} biće zamenjeno starim nazivom
                -	{random} biće zamenjeno slučajno izgenerisanim nazivom
                -	{cDate} biće zamenjeno datumom kreiranja fajla
                
           Primer izraza:  image_{name}_{cDate}
           Novo ime:       image_staroIme_16-07-2019

            
        *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje
           originalnih naziva
           
### 4. Brisanje
       Brisanje svih fajlova i foldera na zadatoj putanji (rekurzivno ili ne) koji su:
       a) Prazni
       b) Kreirani pre zadatog datuma
       c) Nisu korišćeni nakon zadatog datuma

         
Sve operacije bi bile implementirane u programskom jeziku *Go(lang)* uz paralelizaciju (*goroutine*) gde je to moguće (npr. prilikom učitavanja fajlova i foldera i kreiranja stabla).

GUI deo aplikacije iz kog je moguće pozivanje svih navedenih operacija bi bio implementiran uz oslonac na programski jezik *Pharo*. 
Korisnicima će biti omogućen rad korišćenjem odgovarajućih formi, a nešto naprednijim korisnicima biće ponuđeno i tekstualno polje
u koje je moguće unositi komande jednostavnog jezika specifičnog za domen koji će biti razvijen.

Takođe bi bilo moguće i pozivanje operacija(2, 3 i 4) iz komandne linije (implementacija korišćenjem [*Cobra*](https://github.com/spf13/cobra) *Golang* biblioteke).
