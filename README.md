# Analiza i organizacija fajl sistema

### 1. Analiza i vizualizacija fajl sistema
      a) Prikaz u vidu stabla počeviši od odabrane putanje
      b) Bar Chart i Pie Chart prikazi po:
          - Tipovima fajla
          - Zauzeću memorije
          - Datumima kreiranja/korišćenja

### 2. Preimenovanje svih foldera/fajlova na zadatoj putanji
        a) Generisanjem random naziva
        b) Zadavanjem prefiksa/sufiksa koji se dodaju na postojeći naziv
        c) Uklanjanjem zadatog dela naziva ili zamenom tog dela naziva novim
        d) Zadavanjem različitih izraza poput:
            image_[increment_from_1] -> generiše image_1, image_2 ...
            
        *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje originalnih naziva

### 3. Reorganizacija strukture foldera
       Mogućnost organizovanja fajlova na zadatoj putanji (rekurzivno ili ne) tako što bi se odvojili u posebne foldere na osnovu:
         a) Tipa fajla
         b) Veličine fajla
         c) Datuma kreiranja (po mesecima ili godinama) 
         
         *) Informacije o izmenama sačuvati u posebnom fajlu na osnovu kog bi bilo omogućeno vraćanje originalne strukture
         
Sve operacije bi bile implementirane u programskom jeziku *Go(lang)* uz paralelizaciju gde je to moguće (npr. prilikom učitavanja i kreiranja stabla).

Vizualizacije bi bile implementirane uz oslonac na programski jezik *Pharo*. 
