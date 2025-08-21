import requests
import json
import argparse
from datetime import datetime, timezone

def get_usd_rate():
    """Obtiene la tasa del dólar desde DolarAPI"""
    try:
        response = requests.get("https://ve.dolarapi.com/v1/dolares/oficial", timeout=5)
        response.raise_for_status()
        data = response.json()
        
        return {
            'currency': 'USD',
            'price': data.get('promedio', 0),
            'last_updated': data.get('fechaActualizacion', ''),
            'source': data.get('fuente', 'oficial'),
            'name': data.get('nombre', 'Oficial')
        }
    except Exception as e:
        print(f"Error al obtener el dólar: {e}", file=sys.stderr)
        return None

def get_eur_rate():
    """Obtiene la tasa del euro desde una fuente confiable"""
    try:
        # Primero intentamos con DolarAPI (por si acaso)
        try:
            response = requests.get("https://ve.dolarapi.com/v1/dolares/euro", timeout=3)
            if response.status_code == 200:
                data = response.json()
                return {
                    'currency': 'EUR',
                    'price': data.get('promedio', 0),
                    'last_updated': data.get('fechaActualizacion', ''),
                    'source': data.get('fuente', 'oficial'),
                    'name': data.get('nombre', 'Euro')
                }
        except:
            pass  # Si falla, continuamos con el método alternativo

        # Método alternativo: usar una API de cambio de divisas
        response = requests.get("https://api.exchangerate-api.com/v4/latest/EUR")
        response.raise_for_status()
        data = response.json()
        
        # Obtener la tasa EUR a USD
        eur_to_usd = data['rates'].get('USD', 1.0)
        
        # Obtener la tasa USD a VES
        usd_data = get_usd_rate()
        if not usd_data:
            return None
            
        # Calcular EUR a VES: (EUR -> USD) * (USD -> VES)
        eur_to_ves = eur_to_usd * usd_data['price']
        
        return {
            'currency': 'EUR',
            'price': eur_to_ves,
            'last_updated': datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ'),
            'source': 'Exchangerate-API + DolarAPI',
            'name': 'Euro (calculado con tasas reales)'
        }
        
    except Exception as e:
        print(f"Error al obtener el euro: {e}", file=sys.stderr)
        return None

def get_exchange_rate(currency='usd'):
    """Obtiene los datos de tipo de cambio para la moneda especificada"""
    currency = currency.lower()
    
    if currency == 'usd':
        return get_usd_rate()
    elif currency == 'eur':
        return get_eur_rate()
    else:
        print(f"Moneda no soportada: {currency}", file=sys.stderr)
        return None

def get_all_rates():
    """Obtiene las tasas de cambio para USD y EUR"""
    usd_data = get_exchange_rate('usd')
    eur_data = get_exchange_rate('eur')
    
    if not usd_data or not eur_data:
        return None
        
    return {
        'USD': usd_data,
        'EUR': eur_data
    }

def convert_currency(amount, from_currency, to_currency, rates):
    """Convierte entre diferentes monedas usando las tasas proporcionadas"""
    if from_currency == to_currency:
        return amount
        
    # Si estamos convirtiendo a VES
    if to_currency == 'VES':
        return amount * rates[from_currency]['price']
    # Si estamos convirtiendo desde VES
    elif from_currency == 'VES':
        return amount / rates[to_currency]['price']
    # Si estamos convirtiendo entre USD y EUR
    else:
        # Primero a VES, luego a la moneda objetivo
        ves_amount = amount * rates[from_currency]['price']
        return ves_amount / rates[to_currency]['price']

def main():
    # Configurar el parser de argumentos
    parser = argparse.ArgumentParser(description='Obtener tasas de cambio de moneda')
    parser.add_argument('--currency', type=str, default='usd',
                      help='Código de moneda (usd o eur)')
    
    args = parser.parse_args()
    
    # Validar moneda
    currency = args.currency.lower()
    if currency not in ['usd', 'eur']:
        print(json.dumps({
            'error': 'Moneda no soportada. Use "usd" o "eur"',
            'success': False
        }))
        return
    
    # Obtener datos de la moneda solicitada
    data = get_exchange_rate(currency)
    if not data:
        print(json.dumps({
            'error': 'No se pudo obtener la tasa de cambio',
            'success': False
        }), file=sys.stderr)
        return
    
    # Devolver los datos en formato JSON
    print(json.dumps({
        'success': True,
        'data': data
    }, default=str))  # Usamos default=str para manejar objetos datetime

if __name__ == "__main__":
    import sys
    main()