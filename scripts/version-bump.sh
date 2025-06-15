#!/bin/bash
# Script per bump automatico di versione
# Uso: ./version-bump.sh [major|minor|patch] [optional-description]

set -e

# Colori per output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Funzione di help
show_help() {
    echo -e "${YELLOW}❌ Tipo di bump richiesto: major, minor, o patch${NC}"
    echo -e "${BLUE}📖 Uso: ./version-bump.sh [major|minor|patch] [optional-description]${NC}"
    echo ""
    echo -e "${BLUE}📋 Esempi:${NC}"
    echo -e "  ./version-bump.sh patch \"Fix authentication bug\""
    echo -e "  ./version-bump.sh minor \"Add user management feature\""
    echo -e "  ./version-bump.sh major \"Breaking API changes\""
}

# Verifica parametri
if [ -z "$1" ]; then
    show_help
    exit 1
fi

BUMP_TYPE="$1"
DESCRIPTION="$2"

# Verifica che il tipo di bump sia valido
if [[ "$BUMP_TYPE" != "major" && "$BUMP_TYPE" != "minor" && "$BUMP_TYPE" != "patch" ]]; then
    echo -e "${RED}❌ Tipo di bump non valido: $BUMP_TYPE${NC}"
    echo -e "${GREEN}✅ Tipi validi: major, minor, patch${NC}"
    exit 1
fi

echo -e "${CYAN}🔍 Lettura versione corrente...${NC}"

# Legge la versione corrente dal CHANGELOG
if [ ! -f "CHANGELOG.md" ]; then
    echo -e "${RED}❌ File CHANGELOG.md non trovato${NC}"
    exit 1
fi

CURRENT_VERSION=$(grep -E "^## v[0-9]+\.[0-9]+\.[0-9]+" CHANGELOG.md | head -1 | awk '{print $2}')

if [ -z "$CURRENT_VERSION" ]; then
    echo -e "${RED}❌ Impossibile trovare la versione corrente nel CHANGELOG${NC}"
    exit 1
fi

# Rimuove il prefisso 'v' dalla versione
VERSION_NUMBERS="${CURRENT_VERSION#v}"

# Estrae major, minor, patch
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUMBERS"

# Calcola la nuova versione
case "$BUMP_TYPE" in
    "major")
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    "minor")
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    "patch")
        PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="v$MAJOR.$MINOR.$PATCH"

echo -e "${GREEN}📈 Bump da $CURRENT_VERSION a $NEW_VERSION ($BUMP_TYPE)${NC}"

# Se non è stata fornita una descrizione, la chiede
if [ -z "$DESCRIPTION" ]; then
    read -p "📝 Descrizione per questa release: " DESCRIPTION
    if [ -z "$DESCRIPTION" ]; then
        DESCRIPTION="Release $NEW_VERSION"
    fi
fi

# Data corrente in formato ISO
CURRENT_DATE=$(date +%Y-%m-%d)

echo ""
echo -e "${YELLOW}🚀 Creazione release $NEW_VERSION - $DESCRIPTION${NC}"
echo -e "${YELLOW}📅 Data: $CURRENT_DATE${NC}"
echo ""

# Aggiorna il CHANGELOG
echo -e "${CYAN}🔄 Aggiornamento CHANGELOG...${NC}"

# Crea un backup del CHANGELOG
cp CHANGELOG.md CHANGELOG.md.backup

# Crea il nuovo contenuto del CHANGELOG
{
    echo "# Changelog - Automazione e Qualità del Codice"
    echo ""
    echo "## [Unreleased]"
    echo ""
    echo "## $NEW_VERSION - $DESCRIPTION ($CURRENT_DATE)"
    echo ""

    # Copia il contenuto di [Unreleased] sotto la nuova versione
    awk '
    /^## \[Unreleased\]/ { in_unreleased = 1; next }
    /^## v[0-9]/ {
        if (in_unreleased) {
            in_unreleased = 0
            found_first_version = 1
        }
        if (found_first_version) print
        next
    }
    in_unreleased && NF > 0 { print }
    found_first_version { print }
    ' CHANGELOG.md.backup
} > CHANGELOG.md

# Rimuove il backup
rm CHANGELOG.md.backup

echo -e "${GREEN}✅ CHANGELOG aggiornato con la nuova versione${NC}"

# Verifica che Git sia pulito
if [ -n "$(git status --porcelain | grep -v CHANGELOG.md)" ]; then
    echo -e "${YELLOW}⚠️  Ci sono modifiche non committate. Commit prima di continuare.${NC}"
    git status
    read -p "Continuare comunque? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}❌ Operazione annullata${NC}"
        exit 1
    fi
fi

# Crea Git tag
echo -e "${CYAN}🏷️  Creazione Git tag...${NC}"

git add CHANGELOG.md
git commit -m "release: $NEW_VERSION - $DESCRIPTION"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION - $DESCRIPTION"

echo ""
echo -e "${GREEN}✅ Bump di versione completato!${NC}"
echo -e "${GREEN}📋 Nuova versione: $NEW_VERSION${NC}"
echo -e "${GREEN}🏷️  Tag Git creato: $NEW_VERSION${NC}"
echo ""
echo -e "${YELLOW}📤 Per pubblicare:${NC}"
echo -e "${YELLOW}  git push origin main${NC}"
echo -e "${YELLOW}  git push origin $NEW_VERSION${NC}"
