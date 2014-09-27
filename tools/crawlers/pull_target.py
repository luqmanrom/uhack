import sys
import re
import json
import urllib2


def pull_info(company, query):
	search_type = company
	search_query = query
	dpci_arr = {} 
	if (search_type == "walmart"):
		api_key = "26q874epnt7m4uq9cww59sya"
		data = json.load(urllib2.urlopen("http://walmartlabs.api.mashery.com/v1/search?query=" + search_query + "&format=json&apiKey=26q874epnt7m4uq9cww59sya"))
		for key,value in data.iteritems():
			if (key == "items"):
				prod_tuple = []
				for i in value:
					name = str(i.get("name"))
					price = i.get("salePrice")
					try:
						price = float(price)
						price = int(price*100)
						dpci_arr[name] = price 
					except:
						continue
		return dpci_arr
	elif (search_type == "walgreens"):
		api_key = ""
		data = json.load(urllib2.urlopen(""))		
		#for key,value in data.iteritems():

	elif (search_type == "target"):
		api_key = "J5PsS2XGuqCnkdQq0Let6RSfvU7oyPwF"
		data = json.load(urllib2.urlopen("http://api.target.com/v2/products/search?searchTerm=" + search_query + "&sortBy=relevance&pageNumber=1&pageSize=30&key=" + api_key))
		#res = json.dumps(data, sort_keys=True, indent=4, separators=(',', ': '))
		for key,value in data.iteritems():
			if (key == "CatalogEntryView"):
				for i in value:
					dpci = i.get("DPCI")
					if (dpci != None):
						name = i.get("title")
						sdata = json.load(urllib2.urlopen("http://api.target.com/v2/products/"+dpci+"?idType=DPCI&key="+api_key))
						price = sdata.get("CatalogEntryView")[0].get("Offers")[0].get("OfferPrice")[0].get("priceValue")
						try:
							price = int(price)
							dpci_arr[name] = price;
						except:
							continue
		return dpci_arr

def format_info(arr):	
	if arr == []:
		return
	else:
		for key,value in arr.iteritems():
			print key, value 	
			print "\n"


def getTokens(q):
		toks = q.lower().split(' ')
		toks = [getTokens.regex.sub('', tok) for tok in toks]
		return filter(bool, toks)
getTokens.regex = re.compile(r'[^0-9a-z]')

def puller(store, query):
	qu = pull_info(store, query)
	fd = open(store + ".txt", "a")
	for key, value in qu.iteritems():
		fd.write(query+',%s,;%d' % (','.join(getTokens(key)), value) + "\n")
	fd.close()	

if __name__ == "__main__":
	qu = pull_info(sys.argv[1], sys.argv[2])
	fd = open(sys.argv[1] + '.txt', "a")
	for key,value in qu.iteritems():
		fd.write(sys.argv[2]+',%s,;%d' % (','.join(getTokens(key)), value) + "\n")
	fd.close()

	

