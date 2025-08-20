import requests
import json
import argparse

def get_exchange_rate(currency='usd'):
    """Obtiene los datos de tipo de cambio desde la API de pyDolarVenezuela"""
    url = f"https://pydolarve.org/api/v2/tipo-cambio?currency={currency}&format_date=default&rounded_price=true"
    
    try:
        response = requests.get(url)
        response.raise_for_status()
        data = response.json()
        # Asegurarse de que el símbolo sea consistente
        data['currency'] = currency.upper()
        return data
    except requests.exceptions.RequestException as e:
        print(f"Error al obtener los datos: {e}")
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
        }))
        return
    
    # Devolver los datos en formato JSON
    print(json.dumps({
        'success': True,
        'data': data
    }))

if __name__ == "__main__":
    main()