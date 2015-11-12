'use strict';

/* Services */
angular.module('optima')
  .factory('Advertiser', AdvertiserService)
  .factory('AddAdvertiser', AddAdvertiserService);

function AdvertiserService($resource, $auth){
    // this return the json data info of the advertisers
    return $resource('advertisers', {}, {
        query: {
            method:'GET',           
            isArray:true, 
        }
    });
}

function AddAdvertiserService($resource, $auth){
    
    // this save the data of advertisers on the database
    return $resource('addadvertisers/', {}, {
        put: {
            method:'POST',           
            headers: { 'Content-Type': 'application/json' },    
            params : {
                name: '@names',  
                age: '@age',
                sex: '@sex',
                nse: '@nse',
                coverage: '@coverage',
                interest: '@interest',
                category: '@category',
                budget: '@budget',
                objetives: '@objetives',
            },
        }
    });
}
