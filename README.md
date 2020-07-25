# Analiza i organizacija fajlova i foldera

### Autor
    Ana Mijailović, SW13/2016
    mijailovic.sw13.2016@uns.ac.rs

### 1. Analiza i vizualizacija fajlova i foldera
      a) Prikaz u vidu stabla počevši od odabrane putanje
      b) Bar Chart i Pie Chart prikazi po:
          - Tipovima fajla
          - Zauzeću memorije
          - Datumima kreiranja/korišćenja

### 2. Reorganizacija strukture foldera
       Mogućnost organizovanja fajlova na zadatoj putanji (rekurzivno ili ne) tako što bi se odvojili u 
       posebne foldere na osnovu:
         a) Tipa fajla
         b) Veličine fajla
         c) Datuma kreiranja (po mesecima ili godinama) 
         
         *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje 
            originalne strukture
            
### 3. Preimenovanje svih foldera/fajlova na zadatoj putanji
        a) Generisanjem random naziva
        b) Zadavanjem prefiksa/sufiksa koji se dodaju na postojeći naziv
        c) Uklanjanjem zadatog dela naziva ili zamenom tog dela naziva novim
        d) Zadavanjem različitih izraza poput:
            image_[increment_from_1] -> generiše image_1, image_2 ...
            
        *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje
           originalnih naziva
           
### 4. Brisanje
       Brisanje svih fajlova i foldera na zadatoj putanji (rekurzivno ili ne) koji su:
       a) Prazni
       b) Kreirani pre zadatog vremena
       c) Nisu korišćeni zadati vremenski period

         
Sve operacije bi bile implementirane u programskom jeziku *Go(lang)* uz paralelizaciju (*goroutine*) gde je to moguće (npr. prilikom učitavanja fajlova i foldera i kreiranja stabla).

GUI deo aplikacije iz kog je moguće pozivanje svih navedenih operacija bi bio implementiran uz oslonac na programski jezik *Pharo*. 

Takođe bi bilo moguće i pozivanje operacija(2, 3 i 4) iz komandne linije (implementacija korišćenjem [*Cobra*](https://github.com/spf13/cobra) *Golang* biblioteke).
